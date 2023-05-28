package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func String(v string) *string { return &v }

type GithubDB struct {
	Name   string
	Repo   string
	Token  string
	Client *github.Client
}

func NewGithubDB(dbStr string) *GithubDB {
	//Split dbStr to github config.
	settings := strings.Split(dbStr, "/")

	// tokenize db string:
	// name/repo/token
	if len(settings) != 3 {
		log.Println("Github DB String setting failures.")
		return nil
	}

	// Using token to create github client.
	client := createGithubClient(settings[2])

	if client == nil {
		log.Println("Github client init failure.")
	}

	return &GithubDB{
		Name:   settings[0],
		Repo:   settings[1],
		Token:  settings[1],
		Client: client,
	}
}

func (u *GithubDB) Add(user UserFavorite) {
	// Check if user exist.
	if v, num, err := u.getIssue(user.UserId); err != nil {
		//exist.
		log.Println("User:", user.UserId, "exist, num=", num, " v=", v)
		return
	}
	var body string
	if len(user.Favorites) > 0 {
		body = mergeToContent(user.Favorites)
	}
	if err := u.saveIssue(user.UserId, body); err != nil {
		log.Println("saveIssue err:", err)
	}
}

func (u *GithubDB) Get(uid string) (result *UserFavorite, err error) {
	if v, _, err := u.getIssue(uid); err != nil {
		//cannot find.
		log.Println("cannot find any DB, err:", err)
		return nil, err
	} else {
		favs := splitMultiContent(v)

		log.Println("All Fav:", favs)
		return &UserFavorite{
			UserId:    uid,
			Favorites: favs,
		}, nil
	}

}

// ShowAll: Print all result.
func (u *GithubDB) ShowAll() (result []UserFavorite, err error) {
	log.Println("***Get All DB- Not support now.")

	return nil, nil
}

func (u *GithubDB) Update(user *UserFavorite) (err error) {
	title := user.UserId
	content := mergeToContent(user.Favorites)

	if _, num, err := u.getIssue(title); err != nil {
		//Not exist, save new one.
		return u.saveIssue(title, content)
	} else {
		return u.updateIssue(num, title, content)

	}
}

func createGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func splitMultiContent(cnt string) []string {
	return strings.Split(strings.ReplaceAll(cnt, "\r\n", "\n"), "\n")
}

func mergeToContent(links []string) string {
	return strings.Join(links[:], "\n")
}

func (u *GithubDB) saveIssue(title, body string) error {
	input := &github.IssueRequest{
		Title:    String(title),
		Body:     String(body),
		Assignee: String(""),
	}

	_, _, err := u.Client.Issues.Create(context.Background(), u.Name, u.Repo, input)
	if err != nil {
		log.Printf("Issues.Create returned error: %v", err)
		return err
	}
	return nil
}

func (u *GithubDB) getIssue(title string) (string, int, error) {
	ret, _, err := u.Client.Search.Issues(context.Background(), title, nil)
	if err != nil {
		log.Printf("Issues.search returned error: %v", err)
		return "", 0, err
	}

	log.Println("issue ret:", ret)
	for _, v := range ret.Issues {
		log.Println("Issue Num:", v.Number)
		log.Println("Body:", v.Body)
		log.Println("Comments:", v.Comments)
	}
	if len(ret.Issues) == 0 {
		return "", 0, errors.New("NO_Data")
	}
	return *ret.Issues[0].Body, *ret.Issues[0].Number, nil
}

func (u *GithubDB) updateIssue(number int, title string, updatedCnt string) error {
	updateIssue := &github.IssueRequest{
		Title:    String(title),
		Body:     String(updatedCnt),
		Assignee: String(""),
	}
	ret, _, err := u.Client.Issues.Edit(context.Background(), u.Name, u.Repo, number, updateIssue)
	if err != nil {
		fmt.Printf("Issues.edit returned error: %v", err)
		return err
	}

	log.Println("Issue updated:", ret)
	return nil
}
