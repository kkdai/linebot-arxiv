package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/marvin-hansen/arxiv/v1"
)

const (
	SearchByID = "https://export.arxiv.org/api/query?id_list="
	GetNewest  = "http://export.arxiv.org/api/query?search_query=all:electron&start=0&max_results=10&sortBy=submittedDate&sortOrder=descending"
	Random100  = "http://export.arxiv.org/api/query?search_query=all:electron&start=0&max_results=100"
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

	log.Println("Going to:", SearchByID+idStr)
	resp, err := http.Get(SearchByID + idStr)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	var entry arxiv.Feed
	xml.Unmarshal(data, &entry)
	log.Println("Get by ID for Title:", entry.Entry[0].Title)
	return entry.Entry
}

func getNewest10Articles() []*arxiv.Entry {
	log.Println("Going to:", GetNewest)
	resp, err := http.Get(GetNewest)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	var entry arxiv.Feed
	xml.Unmarshal(data, &entry)
	return entry.Entry
}

func getRandom10Articles() []*arxiv.Entry {
	log.Println("Going to:", Random100)
	resp, err := http.Get(Random100)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	var entry arxiv.Feed
	xml.Unmarshal(data, &entry)

	rands := GetRandomIntSet(100, 10)

	var ret []*arxiv.Entry
	for i := 0; i < 10; i++ {
		item := arxiv.Entry{
			Title: entry.Entry[rands[i]].Title,
			ID:    entry.Entry[rands[i]].ID,
			Summary: &arxiv.Text{
				Body: entry.Entry[rands[i]].Summary.Body,
				Type: entry.Entry[rands[i]].Summary.Type},
		}

		ret = append(ret, &item)
	}
	return ret
}

// NormalizeArxivURL takes a URL string and returns a normalized arXiv URL in the format https://arxiv.org/abs/xxxx.xxx.x
func NormalizeArxivURL(link string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	// Check if the host is arxiv.org or huggingface.co
	if parsedURL.Host != "arxiv.org" && parsedURL.Host != "huggingface.co" {
		return "", errors.New("URL does not belong to arxiv.org or huggingface.co")
	}

	// Extract the paper ID using a regular expression
	re := regexp.MustCompile(`(\d{4}\.\d{4,5}(v\d+)?)`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if matches == nil {
		return "", errors.New("URL does not contain a valid arXiv ID")
	}
	paperID := matches[1]

	// Return the normalized arXiv URL
	return "https://arxiv.org/abs/" + paperID, nil
}

// ConvertToPDFURL converts an arXiv URL to its PDF URL
func ConvertToPDFURL(absURL string) (string, error) {
	// First normalize the URL to ensure it's a valid arXiv URL
	normalizedURL, err := NormalizeArxivURL(absURL)
	if err != nil {
		return "", err
	}

	// Extract the paper ID from the normalized URL
	paperID := getIDfromURL(normalizedURL)
	if paperID == "" {
		return "", errors.New("failed to extract paper ID from URL")
	}

	// Return the PDF URL
	return "https://arxiv.org/pdf/" + paperID + ".pdf", nil
}

// IsArxivPDFURL checks if a URL is a valid arXiv PDF URL
func IsArxivPDFURL(link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	// Check if it's arxiv.org and path contains /pdf/
	if parsedURL.Host != "arxiv.org" {
		return false
	}

	// Check if path contains /pdf/ and ends with .pdf
	re := regexp.MustCompile(`/pdf/\d{4}\.\d{4,5}(v\d+)?\.pdf`)
	return re.MatchString(parsedURL.Path)
}
