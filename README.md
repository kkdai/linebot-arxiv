ChatGPT LINE Bot for Arxiv
==============

[![GoDoc](https://godoc.org/github.com/kkdai/linebot-arxiv.svg?status.svg)](https://godoc.org/github.com/kkdai/linebot-arxiv)  ![Go](https://github.com/kkdai/linebot-arxiv/workflows/Go/badge.svg) [![goreportcard.com](https://goreportcard.com/badge/github.com/kkdai/linebot-arxiv)](https://goreportcard.com/report/github.com/kkdai/LineBotTemplate)

Featues
=============

## Menu

<img src="./img/1-menu.jpg" alt="Show" title="Display Menu" width="300" />

Press "menu" to display menu

### Newest

Show the newest article on arXiv.

### Random Articles

Show the random 10 articles on arXiv.

### My Favs

<img src="./img/2-fav.jpg" alt="Fav" title="Show your favorite articles list " width="300" />

Show your favorite articles list.

## Save your arXiv link to Favorite list

<img src="./img/3-save.jpg" alt="Save" title="Save to your fav list" width="300" />

Save to your fav list.

## Query keywords on arXiv

<img src="./img/4-query.jpg" alt="Query" title="query" width="300" />

Query keywords to find articles on arXiv.

How to build your own LINE Bot?
=============

### To obtain a LINE Bot API developer account

Make sure you are registered on the LINE developer console at <https://developers.line.biz/console/> if you want to use a LINE Bot.

Create a new Messaging Channel and get the "Channel Secret" on the "Basic Setting" tab.

Issue a "Channel Access Token" on the "Messaging API" tab.

Open the LINE OA manager from the "Basic Setting" tab and go to the Reply setting on the OA manager. Enable "webhook" there.

### To obtain an Gemini API token

Register for an account on the Google AI Studio website at <https://aistudio.google.com/app/apikey/>.

Once you have an account, you can find your API token in the account settings page.

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
