package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/marvin-hansen/arxiv/v1"
)

// getArxivArticle:
func getArxivArticle(keyword string) []*arxiv.Entry {
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

func getArticleByURL(urlStr string) []*arxiv.Entry {
	idStr := getIDfromURL(urlStr)
	if idStr == "" {
		return nil
	}

	log.Println("Going to:", "https://export.arxiv.org/api/query?id_list="+idStr)
	resp, err := http.Get("https://export.arxiv.org/api/query?id_list=" + idStr)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	log.Println("data=:", data)

	var entry arxiv.Feed
	xml.Unmarshal(data, &entry)

	log.Println("Title:", entry.Entry[0].Title)
	log.Println("Summary:", entry.Entry[0].Summary)
	log.Println("Authors:")
	for _, author := range entry.Entry[0].Author {
		log.Println(" -", author.Name)
	}
	return entry.Entry
}
