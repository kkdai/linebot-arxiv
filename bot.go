package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const PROMPT_AskLocation = `請問以下的文字包含某個地點嗎？ 如果不是的話，請簡單的回覆我 "NO" 即可。 ---- %s ----`
const PROMPT_AskDate = `請問以下的文字包含某段時間嗎？ 如果不是的話，請簡單的回覆我 "NO" 即可。 ---- %s ----`
const PROMPT_AskPeriod = `請問以下的文字包含某段時間嗎？ 如果不是的話，請簡單的回覆我 "NO" 即可。 ---- %s ----`
const PROMPT_AskPlanning = `你是一個旅行社的員工，協助評估顧客的旅遊景點規劃 。現在我即將去: 
地點： 
%s 

期間在: 
%s

總共天數是:
%s

%s 又有什麼特殊的節日？
根據這些節日，有沒有必去的景點規劃? 幫我每一天分開列出。`

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

				// Directly to ChatGPT
				if strings.Contains(message.Text, ":gpt") {
					handleGPT(GPT_Complete, event, message.Text)
				}
			}
		}
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
