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

	bodyPtr := body.String()

	issueReq := &github.IssueRequest{
		Title:  &title,
		Body:   &bodyPtr,
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
	issueList, _, err := githubClient.Issues.ListByRepo(ctx, GithubUsername, GithubRepoName, &github.IssueListByRepoOptions{})
	if err != nil {
		return nil, err
	}

	var existingIssue github.Issue

	for _, issue := range issueList {
		if strings.Contains(*issue.Title, sessionId) {
			existingIssue = *issue
		}
	}

	log.Println("Existing issue", existingIssue)

	return &existingIssue, err
}

func commentOnExistingIssue(ctx context.Context, issue *github.Issue, sessionUrl, pageUrl, noteText, author string) error {
	log.Printf("issue struct %+v\\n", issue)

	issueBody := &IssueBody{
		NoteText:   noteText,
		SessionUrl: sessionUrl,
		PageUrl:    pageUrl,
		Author:     author,
		IsComment:  true,
	}

	issueBodyPtr := issueBody.String()

	user, _, err := githubClient.Users.Get(ctx, "")
	if err != nil {
		log.Println("error fetching user", err)
		return err
	}

	log.Println("repo owner", github.Stringify(issue.Repository.Owner.Name))
	log.Println("repo string", github.Stringify(issue.Repository.Name))
	log.Println("issue number", issue.GetNumber())

	_, _, err = githubClient.Issues.CreateComment(ctx, github.Stringify(issue.Repository.Owner.Name), github.Stringify(issue.Repository.Name), issue.GetNumber(), &github.IssueComment{
		ID:                nil,
		NodeID:            nil,
		Body:              github.String(issueBodyPtr),
		User:              user,
		Reactions:         nil,
		CreatedAt:         nil,
		UpdatedAt:         nil,
		AuthorAssociation: nil,
		URL:               nil,
		HTMLURL:           nil,
		IssueURL:          nil,
	})
	if err != nil {
		log.Println("error creating comment on issue", err)
		return err
	}
	log.Println("issue request func returned")
	return nil
}
