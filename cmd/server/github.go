package main

import (
	"bytes"
	"context"
	"log"
	"strings"
	"text/template"

	"github.com/google/go-github/v35/github"
)

const (
	GithubUsername = "KingsleyBawuah"
	GithubRepoName = "MovieSearch"
)

type IssueBody struct {
	NoteText   string
	SessionUrl string
	PageUrl    string
	Author     string
	IsComment  bool
}

func (i IssueBody) String() string {
	tmpl, err := template.New("issueBody").Parse(`
### Note Text: 
	{{.NoteText}}

### Relevant Links: 

Session Url:
- {{.SessionUrl}}

Page where error occurred:
- {{.PageUrl}}



_{{if .IsComment}}Comment {{else}}Issue {{end}}created automatically from a note in Fullstory using the #issue command by the author: {{.Author}}_
`)

	if err != nil {
		log.Panic("error creating issue body template", err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, i)

	return buf.String()
}

func createGithubIssue(ctx context.Context, title, sessionUrl, pageUrl, noteText, author string) error {
	labels := []string{"bug", "auto-generated", "bot", "fullstory"}

	body := &IssueBody{
		NoteText:   noteText,
		SessionUrl: sessionUrl,
		PageUrl:    pageUrl,
		Author:     author,
		IsComment:  false,
	}

	bodyString := body.String()

	issueReq := &github.IssueRequest{
		Title:  &title,
		Body:   &bodyString,
		Labels: &labels,
	}

	_, _, err := githubClient.Issues.Create(ctx, GithubUsername, GithubRepoName, issueReq)
	if err != nil {
		return err
	}

	return nil
}

// Determine if a issue already exists for a specified session ID, if so return it.
func inquireExistingIssue(ctx context.Context, sessionId string) (*github.Issue, error) {
	issueList, _, err := githubClient.Issues.ListByRepo(ctx, GithubUsername, GithubRepoName, &github.IssueListByRepoOptions{
		State: "open",
	})
	if err != nil {
		return nil, err
	}

	var existingIssue *github.Issue

	for _, issue := range issueList {
		if strings.Contains(*issue.Title, sessionId) {
			existingIssue = issue
		}
	}
	return existingIssue, err
}

func commentOnExistingIssue(ctx context.Context, issue *github.Issue, sessionUrl, pageUrl, noteText, author string) error {
	issueBody := &IssueBody{
		NoteText:   noteText,
		SessionUrl: sessionUrl,
		PageUrl:    pageUrl,
		Author:     author,
		IsComment:  true,
	}

	issueBodyString := issueBody.String()

	_, _, err := githubClient.Issues.CreateComment(ctx, GithubUsername, GithubRepoName, issue.GetNumber(), &github.IssueComment{
		Body: github.String(issueBodyString),
	})
	if err != nil {
		log.Println("error creating comment on issue", err)
		return err
	}
	return nil
}
