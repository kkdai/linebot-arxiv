ChatGPT LINE Bot for Arxiv
==============

[![GoDoc](https://godoc.org/github.com/kkdai/linebot-arxiv.svg?status.svg)](https://godoc.org/github.com/kkdai/linebot-arxiv)  ![Go](https://github.com/kkdai/linebot-arxiv/workflows/Go/badge.svg) [![goreportcard.com](https://goreportcard.com/badge/github.com/kkdai/linebot-arxiv)](https://goreportcard.com/report/github.com/kkdai/LineBotTemplate)

How to use this
=============

### To obtain a LINE Bot API developer account

Make sure you are registered on the LINE developer console at <https://developers.line.biz/console/> if you want to use a LINE Bot.

Create a new Messaging Channel and get the "Channel Secret" on the "Basic Setting" tab.

Issue a "Channel Access Token" on the "Messaging API" tab.

Open the LINE OA manager from the "Basic Setting" tab and go to the Reply setting on the OA manager. Enable "webhook" there.

### To obtain an OpenAI API token

Register for an account on the OpenAI website at <https://openai.com/api/>.

Once you have an account, you can find your API token in the account settings page.

If you want to use the OpenAI API for development, you can find more information and instructions in the API documentation page.

Please note that the OpenAI API is only available to users who meet certain criteria. You can find more information about the usage conditions and limitations of the API in the API documentation page.

### Deploy this on Web Platform

You can choose [Heroku](https://www.heroku.com/) or [Render](http://render.com/)

#### Deploy this on Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

- Input `Channel Secret` and `Channel Access Token` and `ChatGptToken`.

#### Deploy this on Rener

[![Deploy to Render](http://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)

- Input `Channel Secret` and `Channel Access Token` and `ChatGptToken`.

License
---------------

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
