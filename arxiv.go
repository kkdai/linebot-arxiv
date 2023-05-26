package main

import (
	"log"

	"github.com/jtracks/go-arciv/arciv"
)

// getArxivArticle:
func getArxivArticle(keyword string) []arciv.Entry {
	result, _ := arciv.Search(
		arciv.SimpleQuery{
			Search:     keyword,
			MaxResults: 5,
		})

	for i, e := range result.Entries {
		log.Printf("Result %v: %v\n %v", i+1, e.Title, e.Links)
	}
	return result.Entries
}

//
