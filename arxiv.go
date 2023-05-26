package main

import (
	"fmt"

	"github.com/jtracks/go-arciv/arciv"
)

// goArxivArticle:
func goArxivArticle(keywords string) error {
	result, _ := arciv.Search(
		arciv.SimpleQuery{
			Search:     "electron",
			MaxResults: 5,
		})

	for i, e := range result.Entries {
		fmt.Printf("Result %v: %v\n %v", i+1, e.Title, e.Summary)
	}
	return nil
}
