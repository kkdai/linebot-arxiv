package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/orijtech/arxiv/v1"
	"golang.org/x/tools/blog/atom"
)

// getArxivArticle:
func getArxivArticle(keyword string) []*atom.Entry {
	resChan, _, err := arxiv.Search(context.Background(), &arxiv.Query{
		Terms:         keyword,
		MaxPageNumber: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for resPage := range resChan {
		if err := resPage.Err; err != nil {
			log.Printf("#%d err: %v", resPage.PageNumber, err)
			continue
		}

		log.Printf("#%d\n", resPage.PageNumber)
		feed := resPage.Feed
		log.Printf("\tTitle: %s\n\tID: %s\n\tAuthor: %#v\n\tUpdated: %#v\n", feed.Title, feed.ID, feed.Author, feed.Updated)
		return feed.Entry
	}
	return nil
}

func getIDfromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Remove the leading slash and "abs/" from the path
	return strings.TrimPrefix(u.Path, "/abs/")
}
