package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Intent struct {
	Keywords        string `json:"keywords"`
	NumberOfArticle int    `json:"numberOfArticle"`
}

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

				handleArxivSearch(event, message.Text)

			}
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
	return &intent
}

// handleArxivSearch:
func handleArxivSearch(event *linebot.Event, msg string) {

	// result := getArxivArticle(msg)
	reply := gptCompleteContext(msg)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		log.Print(err)
	}

}

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
