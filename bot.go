package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/kkdai/linebot-arxiv/models"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/marvin-hansen/arxiv/v1"
)

// Postback Actions
const (
	ActionOpenDetail      string = "DetailArticle"
	ActionTransArticle    string = "TransArticle"
	ActionBookmarkArticle string = "BookmarkArticle"
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
				if strings.Contains(message.Text, "arxiv.org/abs") {
					values := url.Values{}
					values.Set("user_id", event.Source.UserID)
					values.Set("url", message.Text)
					values.Set("extra", "gpt")
					actionBookmarkArticle(event, values)
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

// parseIntent:
func parseIntent(msg string) *Intent {
	gpt1 := fmt.Sprintf(PROMPT_GetIntent, msg)
	reply := gptCompleteContext(gpt1)

	var intent Intent
	// Unmarshal the JSON data into the struct
	err := json.Unmarshal([]byte(reply), &intent)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	log.Println(" Intent:=", intent)
	return &intent
}

// handleArxivSearch:
func handleArxivSearch(event *linebot.Event, msg string) {
	results := getArxivArticle(msg)

	template := getCarouseTemplate(event.Source.UserID, results)
	if template != nil {
		sendCarouselMessage(event, template, "Paper Result")
	}
}

// handleGPT:
func handleGPT(action GPT_ACTIONS, event *linebot.Event, message string) {
	switch action {
	case GPT_Complete:
		reply := gptCompleteContext(message)
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
			log.Print(err)
		}
	case GPT_Draw:
		if reply, err := gptImageCreate(message); err != nil {
			if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無法正確顯示圖形.")).Do(); err != nil {
				log.Print(err)
			}
		} else {
			if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("根據你的提示，畫出以下圖片："), linebot.NewImageMessage(reply, reply)).Do(); err != nil {
				log.Print(err)
			}
		}
	}
}

func isGroupEvent(event *linebot.Event) bool {
	return event.Source.GroupID != "" || event.Source.RoomID != ""
}

func getGroupID(event *linebot.Event) string {
	if event.Source.GroupID != "" {
		return event.Source.GroupID
	} else if event.Source.RoomID != "" {
		return event.Source.RoomID
	}

	return ""
}

func getCarouseTemplate(userId string, records []*arxiv.Entry) (template *linebot.CarouselTemplate) {
	if len(records) == 0 {
		log.Println("err1")
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
	content := fmt.Sprintf("論文： %s \n作者: \n %s \n摘要: \n %s \n論文網址: \n%s", result[0].Title, authors, result[0].Summary.Body, result[0].Link[1].Href)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
		log.Println(err)
	}
}

func actionGPTTranslate(event *linebot.Event, values url.Values) {
	url := values.Get("url")
	log.Println("actionGPTTranslate: url=", url)
	result := getArticleByURL(url)
	gptRet := gptCompleteContext(fmt.Sprintf(`幫我將以下內容做中文摘要: ---\n %s---"`, result[0].Summary.Body))
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(gptRet)).Do(); err != nil {
		log.Println(err)
	}
}

func actionBookmarkArticle(event *linebot.Event, values url.Values) {
	newFavoriteArticle := values.Get("url")
	uid := values.Get("user_id")
	extraAct := values.Get("extra")
	var toggleMessage = "已新增至最愛"
	newUser := models.UserFavorite{
		UserId:    uid,
		Favorites: []string{newFavoriteArticle},
	}
	if record, err := DB.Get(uid); err != nil {
		log.Println("User data is not created, create a new one")
		DB.Add(newUser)
		log.Println(newFavoriteArticle, "Add user/fav")
	} else {
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

	var gptRet string
	ret := fmt.Sprintf("文章: \n%s \n%s", newFavoriteArticle, toggleMessage)
	if strings.Compare(extraAct, "gpt") == 0 {
		result := getArticleByURL(newFavoriteArticle)
		gptRet = gptCompleteContext(fmt.Sprintf(`幫我將以下內容做中文摘要: ---\n %s---"`, result[0].Summary.Body))
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret), linebot.NewTextMessage(gptRet)).Do(); err != nil {
			log.Println(err)
		}
	} else {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret)).Do(); err != nil {
			log.Println(err)
		}
	}
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}

// InArray: Check if string item is in array
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

// RemoveStringItem: Remove string item from slice
func RemoveStringItem(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
