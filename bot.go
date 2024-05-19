package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/kkdai/favdb"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/marvin-hansen/arxiv/v1"
)

// Postback Actions
const (
	ActionOpenDetail      string = "DetailArticle"
	ActionTransArticle    string = "TransArticle"
	ActionBookmarkArticle string = "BookmarkArticle"
	ActionHelp            string = "Menu"
	ActonShowFav          string = "MyFavs"
	ActionNewest          string = "Newest"
	ActionRandom          string = "Random"
)

type Intent struct {
	Keywords        string `json:"keywords"`
	NumberOfArticle int    `json:"numberOfArticle"`
}

const Image_URL = "https://github.com/kkdai/linebot-arxiv/blob/f9ca955ff9392f5af4d27e617e5f71fc97c8f60e/img/paper.png?raw=true"
const PROMPT_GetIntent = `幫我把以下文字，拆成 JSON 回覆。 
"%s"
---
{   
keywords: ""
numberOfArticle: 0
}
---`

const PROMPT_Summarization = `幫我將以下內容做中文摘要: ---
 %s
 ---`

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				if isGroupEvent(event) {
					// 如果聊天機器人在群組中，不回覆訊息。
					return
				}
				if refineURL, err := NormalizeArxivURL(message.Text); err == nil {
					values := url.Values{}
					values.Set("user_id", event.Source.UserID)
					values.Set("url", refineURL)
					values.Set("extra", "gpt")
					actionBookmarkArticle(event, values)
					return
				} else if strings.EqualFold(message.Text, "menu") {
					template := getMenuButtonTemplate(event, "論文收集")
					sendCarouselMessage(event, template, "我能為您做什麼？")
					return
				}

				handleArxivSearch(event, message.Text)
			}
		} else if event.Type == linebot.EventTypePostback {
			log.Println("got a postback event")
			log.Println(event.Postback.Data)
			postbackHandler(event)

		}
	}
}

// handleArxivSearch:
func handleArxivSearch(event *linebot.Event, msg string) {
	results := getArxivArticle(msg)

	template := getCarouseTemplate(event.Source.UserID, results)
	if template != nil {
		sendCarouselMessage(event, template, "Paper Result")
	}
}

func getCarouseTemplate(userId string, records []*arxiv.Entry) (template *linebot.CarouselTemplate) {
	if len(records) == 0 {
		log.Println("err: Empty articles.")
		return nil
	}

	var checkList []string
	if record, err := DB.Get(userId); err == nil {
		checkList = record.Favorites
	}

	log.Println("all items:", checkList)

	columnList := []*linebot.CarouselColumn{}
	for _, result := range records {
		var saveTogle string
		if exist, _ := InArray(result.ID, checkList); !exist {
			saveTogle = "儲存文章"
		} else {
			saveTogle = "移除儲存"
		}
		detailData := fmt.Sprintf("action=%s&url=%s&user_id=%s", ActionOpenDetail, result.ID, userId)
		transData := fmt.Sprintf("action=%s&url=%s&user_id=%s", ActionTransArticle, result.ID, userId)
		SaveData := fmt.Sprintf("action=%s&url=%s&user_id=%s", ActionBookmarkArticle, result.ID, userId)
		tmpColumn := linebot.NewCarouselColumn(
			Image_URL,
			truncateString(result.Title, 35)+"..",
			truncateString(result.Summary.Body, 55)+"..",
			//			linebot.NewURIAction("打開網址", result.ID),
			linebot.NewPostbackAction("知道更多", detailData, "", "", "", ""),
			linebot.NewPostbackAction("翻譯摘要(比較久)", transData, "", "", "", ""),
			linebot.NewPostbackAction(saveTogle, SaveData, "", "", "", ""),
		)
		columnList = append(columnList, tmpColumn)
	}
	template = linebot.NewCarouselTemplate(columnList...)
	return template
}

func sendCarouselMessage(event *linebot.Event, template *linebot.CarouselTemplate, altText string) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage(altText, template)).Do(); err != nil {
		log.Println(err)
	}
}

func postbackHandler(event *linebot.Event) {
	m, _ := url.ParseQuery(event.Postback.Data)
	action := m.Get("action")
	log.Println("Action = ", action)
	actionHandler(event, action, m)
}

func actionHandler(event *linebot.Event, action string, values url.Values) {
	switch action {
	case ActionOpenDetail:
		log.Println("ActionOpenDetail:", values)
		actionGetDetail(event, values)
	case ActionTransArticle:
		log.Println("ActionTransArticle:", values)
		actionGPTTranslate(event, values)
	case ActionBookmarkArticle:
		log.Println("ActionBookmarkArticle:", values)
		actionBookmarkArticle(event, values)
		log.Println("Show all article:....")
		DB.ShowAll()
	case ActonShowFav:
		log.Println("ActonShowFav:", values)
		actionShowFavorite(event, values)
	case ActionNewest:
		log.Println("ActionNewest:", values)
		actionNewest(event, values)
	case ActionRandom:
		log.Println("ActionRandom:", values)
		actionRandom(event, values)
	default:
		log.Println("Unimplement action handler", action)
	}
}

func actionGetDetail(event *linebot.Event, values url.Values) {
	url := values.Get("url")
	log.Println("actionGPTTranslate: url=", url)
	result := getArticleByURL(url)
	authors := ""
	for _, a := range result[0].Author {
		authors = fmt.Sprintf("%s\n%s", authors, a.Name)
	}
	content := fmt.Sprintf("論文： %s \n作者: \n %s \n摘要: \n %s \n論文網址: \n%s \nPDF: \n%s", result[0].Title, authors, result[0].Summary.Body, result[0].ID, result[0].Link[1].Href)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
		log.Println(err)
	}
}

// actionGPTTranslate: Translate article by GPT-3.
func actionGPTTranslate(event *linebot.Event, values url.Values) {
	url := values.Get("url")
	log.Println("actionGPTTranslate: url=", url)
	result := getArticleByURL(url)
	sumResult, err := GeminiChat(fmt.Sprintf(PROMPT_Summarization, result[0].Summary.Body))
	if err != nil {
		log.Println("Error:", err)
		errString := fmt.Sprintf("Error: %s", err)
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(errString)).Do()
		return
	}
	//Doing url handle if it in gpt summarization.
	sumResult = AddLineBreaksAroundURLs(sumResult)

	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(sumResult)).Do(); err != nil {
		log.Println(err)
	}
}

// actionBookmarkArticle: Add or remove article from favorite list.
func actionBookmarkArticle(event *linebot.Event, values url.Values) {
	newFavoriteArticle := values.Get("url")
	uid := values.Get("user_id")
	extraAct := values.Get("extra")
	var toggleMessage = "已新增至最愛"
	newUser := favdb.UserFavorite{
		UserId:    uid,
		Favorites: []string{newFavoriteArticle},
	}
	if record, err := DB.Get(uid); err != nil {
		log.Println("User data is not created, create a new one")
		DB.Add(newUser)
		log.Println(newFavoriteArticle, "Add user/fav")
	} else {
		// from link to chatbot. skip removed only show summary.
		if strings.Compare(extraAct, "gpt") != 0 {
			log.Println("Record found, update it", record)
			oldRecords := record.Favorites

			if exist, idx := InArray(newFavoriteArticle, oldRecords); exist == true {
				log.Println(newFavoriteArticle, "Del fav")
				oldRecords = RemoveStringItem(oldRecords, idx)
				toggleMessage = "已從最愛中移除"
			} else {
				log.Println(newFavoriteArticle, "Add fav")
				oldRecords = append(oldRecords, newFavoriteArticle)
			}
			record.Favorites = oldRecords
			DB.Update(record)
		}
	}

	ret := fmt.Sprintf("文章: \n%s \n%s", newFavoriteArticle, toggleMessage)
	if strings.Compare(extraAct, "gpt") == 0 {
		result := getArticleByURL(newFavoriteArticle)
		sumResult, err := GeminiChat(fmt.Sprintf(PROMPT_Summarization, result[0].Summary.Body))
		if err != nil {
			log.Println("Error:", err)
			errString := fmt.Sprintf("Error: %s", err)
			bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(errString)).Do()
			return
		}

		log.Println("Gemini response:", sumResult)
		dataShowFav := fmt.Sprintf("action=%s&user_id=%s&page=0", ActonShowFav, event.Source.UserID)
		qrBookmark := linebot.NewQuickReplyItems(linebot.NewQuickReplyButton("", linebot.NewPostbackAction("列出 My Fav", dataShowFav, "", "", "", "")))
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(sumResult), linebot.NewTextMessage("論文位置在："+newFavoriteArticle+"\n 論文已經儲存。").WithQuickReplies(qrBookmark)).Do(); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("normal response:", ret)
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret)).Do(); err != nil {
			log.Println(err)
		}
	}
}

// actionShowFavorite: Show favorite article list.
func actionShowFavorite(event *linebot.Event, values url.Values) {
	log.Println("actionShowFavorite call")
	columnCount := 9
	userId := values.Get("user_id")

	if currentPage, err := strconv.Atoi(values.Get("page")); err != nil {
		log.Println("Unable to parse parameters", values)
	} else {
		userData, _ := DB.Get(userId)

		// No userData or user has empty Fav, return!
		if userData == nil || (userData != nil && len(userData.Favorites) == 0) {
			empStr := "你沒有收藏任何文章，快來加入吧。"
			// Fav == 0, skip it.
			empColumn := linebot.NewCarouselColumn(
				Image_URL,
				"沒有論文",
				empStr,
				linebot.NewMessageAction(ActionHelp, ActionHelp),
			)
			emptyResult := linebot.NewCarouselTemplate(empColumn)
			sendCarouselMessage(event, emptyResult, empStr)
		}

		startIdx := currentPage * columnCount
		endIdx := startIdx + columnCount
		lastPage := false

		// reverse slice
		for i := len(userData.Favorites)/2 - 1; i >= 0; i-- {
			opp := len(userData.Favorites) - 1 - i
			userData.Favorites[i], userData.Favorites[opp] = userData.Favorites[opp], userData.Favorites[i]
		}

		if endIdx > len(userData.Favorites)-1 || startIdx > endIdx {
			endIdx = len(userData.Favorites)
			lastPage = true
		}

		var favDocuments []*arxiv.Entry
		favs := userData.Favorites[startIdx:endIdx]
		log.Println(favs)
		for i := startIdx; i < endIdx; i++ {
			url := userData.Favorites[i]
			tmpRecord := getArticleByURL(url)
			favDocuments = append(favDocuments, tmpRecord[0])
		}

		// append next page column
		previousPage := currentPage - 1
		if previousPage < 0 {
			previousPage = 0
		}
		nextPage := currentPage + 1
		previousData := fmt.Sprintf("action=%s&page=%d&user_id=%s", ActonShowFav, previousPage, userId)
		nextData := fmt.Sprintf("action=%s&page=%d&user_id=%s", ActonShowFav, nextPage, userId)
		previousText := fmt.Sprintf("上一頁 %d", previousPage)
		nextText := fmt.Sprintf("下一頁 %d", nextPage)
		if lastPage == true {
			nextData = "--"
			nextText = "--"
		}

		tmpColumn := linebot.NewCarouselColumn(
			Image_URL,
			"沒有論文",
			"繼續看？",
			linebot.NewMessageAction(ActionHelp, ActionHelp),
			linebot.NewPostbackAction(previousText, previousData, "", "", "", ""),
			linebot.NewPostbackAction(nextText, nextData, "", "", "", ""),
		)

		template := getCarouseTemplate(event.Source.UserID, favDocuments)
		template.Columns = append(template.Columns, tmpColumn)
		sendCarouselMessage(event, template, "收藏的論文已送達")
	}
}

// actionNewest: Show newest 10 articles.
func actionNewest(event *linebot.Event, values url.Values) {
	results := getNewest10Articles()
	template := getCarouseTemplate(event.Source.UserID, results)
	if template != nil {
		sendCarouselMessage(event, template, "Paper Result")
	}
}

// actionRandom: Show random 10 articles.
func actionRandom(event *linebot.Event, values url.Values) {
	results := getRandom10Articles()
	template := getCarouseTemplate(event.Source.UserID, results)
	if template != nil {
		sendCarouselMessage(event, template, "Paper Result")
	}
}

// getMenuButtonTemplate: Get menu button template.
func getMenuButtonTemplate(event *linebot.Event, title string) (template *linebot.CarouselTemplate) {
	columnList := []*linebot.CarouselColumn{}
	dataNewlest := fmt.Sprintf("action=%s&page=0", ActionNewest)
	dataRandom := fmt.Sprintf("action=%s", ActionRandom)
	dataShowFav := fmt.Sprintf("action=%s&user_id=%s&page=0", ActonShowFav, event.Source.UserID)

	menu1 := linebot.NewCarouselColumn(
		Image_URL,
		title,
		"你可以試試看以下選項，或直接輸入關鍵字查詢",
		linebot.NewPostbackAction(ActionNewest, dataNewlest, "", "", "", ""),
		linebot.NewPostbackAction(ActionRandom, dataRandom, "", "", "", ""),
		linebot.NewPostbackAction(ActonShowFav, dataShowFav, "", "", "", ""),
	)
	columnList = append(columnList, menu1)
	template = linebot.NewCarouselTemplate(columnList...)
	return template
}
