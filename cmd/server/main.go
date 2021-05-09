package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

var githubClient *github.Client

// https://developer.fullstory.com/note-created
type FSNoteCreatedReqBody struct {
	EventName string `json:"eventName"`
	Version   int    `json:"version"`
	Data      struct {
		ID         string    `json:"id"`
		Created    time.Time `json:"created"`
		Author     string    `json:"author"`
		Text       string    `json:"text"`
		SessionUrl string    `json:"sessionUrl"`
		UserUrl    string    `json:"userUrl"`
		ShareLink  string    `json:"shareLink"`
		PageInfo   struct {
			PageUrl    string `json:"pageUrl"`
			Ipaddress  string `json:"ipAddress"`
			Useragent  string `json:"userAgent"`
			Referrer   string `json:"referrer"`
			Country    string `json:"country"`
			PageHeight int    `json:"pageHeight"`
			PageWidth  int    `json:"pageWidth"`
		} `json:"pageInfo"`
		NotedTime time.Time `json:"notedTime"`
	} `json:"data"`
}

func main() {
	// Set up Github Client
	ctx := context.Background()

	githubToken, err := envMust("GITHUB_TOKEN")
	if err != nil {
		os.Exit(1)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubClient = github.NewClient(tc)

	// Set up the HTTP Request Handler.
	port, err := envMust("PORT")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/handleNote", handleNote)

	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

// containsIssueCmd determines if the note text contains the string #issue which indicates this note should create an issue against the site repo.
func containsIssueCmd(text string) bool {
	matchIssueCmd := `\W(\#(issue)+\b)`

	match, err := regexp.MatchString(matchIssueCmd, text)
	if err != nil {
		log.Panicln("error determining if note text contains issue command", err)
	}

	return match
}

func createGithubIssue(ctx context.Context, title, sessionUrl, noteText, author string) error {
	labels := []string{"bug", "auto-generated", "bot", "fullstory"}

	body := fmt.Sprintf(`
### Note Text: 
	%s

### Link to session: 
	%s

_Issue created automatically from a note in Fullstory using the #issue command by the author: %s_
`, noteText, sessionUrl, author)

	issueReq := &github.IssueRequest{
		Title:     &title,
		Body:      &body,
		Labels:    &labels,
		Assignee:  nil,
		Assignees: nil,
	}

	_, _, err := githubClient.Issues.Create(ctx, "KingsleyBawuah", "MovieSearch", issueReq)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Make this function less busy.
func handleNote(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	switch req.Method {
	case "POST":
		decoder := json.NewDecoder(req.Body)
		var body FSNoteCreatedReqBody
		err := decoder.Decode(&body)
		if err != nil {
			log.Panicln("error decoding request body", err)
		}

		if containsIssueCmd(body.Data.Text) {
			log.Println("True clause, contains #issue")
			// Create the github issue.
			if err := createGithubIssue(ctx, fmt.Sprintf("Error in session %s", body.Data.ID), body.Data.ShareLink, body.Data.Text, body.Data.Author); err != nil {
				log.Panicln("error creating github issue", err)
			}
			// TODO: Be deliberate about the response codes returned.
			fmt.Fprintf(w, "Yo that's a cmd")
		} else {
			log.Println("False clause doesn't contain #issue")
			fmt.Fprintf(w, body.EventName)
		}

	default:
		fmt.Fprintf(w, "This application only supports POST requests, please try again :)\n")
	}

}

func envMust(envVar string) (string, error) {
	if v := os.Getenv(envVar); v == "" {
		return "", errors.New(fmt.Sprintf("error Enviorment Variable %s cannot be empty", envVar))
	} else {
		return v, nil
	}
}
