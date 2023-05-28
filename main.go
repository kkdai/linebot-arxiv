// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kkdai/linebot-arxiv/models"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	gpt3 "github.com/sashabaranov/go-openai"
)

var bot *linebot.Client
var client *gpt3.Client

var DB models.UserFavData

type GPT_ACTIONS int

const (
	GPT_Complete GPT_ACTIONS = 0
	GPT_Draw     GPT_ACTIONS = 1
	GPT_Whister  GPT_ACTIONS = 2
)

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)

	port := os.Getenv("PORT")
	apiKey := os.Getenv("ChatGptToken")

	if apiKey != "" {
		client = gpt3.NewClient(apiKey)
	}

	gitUrl := os.Getenv("GITHUB_URL")
	if gitUrl != "" {
		// Use Github Issue as DB.
		DB = models.NewGithubDB(gitUrl)
	} else {
		DB = models.NewMemDB()
	}

	http.HandleFunc("/callback", callbackHandler)
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
