module github.com/kkdai/linebot-arxiv

// +heroku goVersion go1.20
go 1.20

require github.com/line/line-bot-sdk-go/v7 v7.19.0

require (
	github.com/marvin-hansen/arxiv v0.1.2
	github.com/sashabaranov/go-openai v1.5.0
)

require (
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/orijtech/otils v0.0.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/net v0.10.0 // indirect
)
