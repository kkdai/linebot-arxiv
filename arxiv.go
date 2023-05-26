package main

import (
	"fmt"
	"log"

	"github.com/jtracks/go-arciv/arciv"
)

// getArxivArticle:
func getArxivArticle(keyword string) string {
	result, _ := arciv.Search(
		arciv.SimpleQuery{
			Search:     keyword,
			MaxResults: 5,
		})

	ret := "你找到論文結果如下："

	for i, e := range result.Entries {
		log.Printf("Result %v: %v\n %v", i+1, e.Title, e.Links)
		ret = ret + fmt.Sprintf("\n 標題: %s \n 摘要: %s \n link: %s \n\n\n", e.Title, e.Summary, e.Links[0])
	}
	return ret
}

//
