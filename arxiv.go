package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

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

func getIDfromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Remove the leading slash and "abs/" from the path
	return strings.TrimPrefix(u.Path, "/abs/")
}

func getDetailFromID(urlStr string) *arciv.Entry {
	id := getIDfromURL(urlStr)
	if id == "" {
		return nil
	}
	return nil
}

/*
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Link    struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
	} `xml:"link"`
	Title struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"title"`
	ID           string `xml:"id"`
	Updated      string `xml:"updated"`
	TotalResults struct {
		Text       string `xml:",chardata"`
		Opensearch string `xml:"opensearch,attr"`
	} `xml:"totalResults"`
	StartIndex struct {
		Text       string `xml:",chardata"`
		Opensearch string `xml:"opensearch,attr"`
	} `xml:"startIndex"`
	ItemsPerPage struct {
		Text       string `xml:",chardata"`
		Opensearch string `xml:"opensearch,attr"`
	} `xml:"itemsPerPage"`
	Entry struct {
		Text      string `xml:",chardata"`
		ID        string `xml:"id"`
		Updated   string `xml:"updated"`
		Published string `xml:"published"`
		Title     string `xml:"title"`
		Summary   string `xml:"summary"`
		Author    []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
		} `xml:"author"`
		Comment struct {
			Text  string `xml:",chardata"`
			Arxiv string `xml:"arxiv,attr"`
		} `xml:"comment"`
		Link []struct {
			Text  string `xml:",chardata"`
			Href  string `xml:"href,attr"`
			Rel   string `xml:"rel,attr"`
			Type  string `xml:"type,attr"`
			Title string `xml:"title,attr"`
		} `xml:"link"`
		PrimaryCategory struct {
			Text   string `xml:",chardata"`
			Arxiv  string `xml:"arxiv,attr"`
			Term   string `xml:"term,attr"`
			Scheme string `xml:"scheme,attr"`
		} `xml:"primary_category"`
		Category []struct {
			Text   string `xml:",chardata"`
			Term   string `xml:"term,attr"`
			Scheme string `xml:"scheme,attr"`
		} `xml:"category"`
	} `xml:"entry"`
}

func main() {
	resp, err := http.Get("https://export.arxiv.org/api/query?id_list=2305.12720v1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var entry Feed
	xml.Unmarshal(data, &entry)

	fmt.Println("Title:", entry.Entry.Title)
	fmt.Println("Summary:", entry.Entry.Summary)
	fmt.Println("Authors:")
	for _, author := range entry.Entry.Author {
		fmt.Println(" -", author.Name)
	}
	fmt.Println("PDF Link:")
	for _, link := range entry.Entry.Link {
		if link.Title == "pdf" {
			fmt.Println(" -", link.Href)
		}
	}
}
*/
