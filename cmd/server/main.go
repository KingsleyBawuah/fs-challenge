package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

var githubClient *github.Client

func main() {
	ctx := context.Background()

	// Set up Github Client
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
	mux.HandleFunc("/handleNote", noteRequestHandler)

	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func noteRequestHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	switch req.Method {
	case "POST":
		if err := handleNote(ctx, req.Body); err != nil {
			log.Panicln("error handling note request", err)
		}
		w.WriteHeader(http.StatusOK)
	default:
		fmt.Fprintf(w, "This application only supports POST requests, please try again :)\n")
	}

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

func handleNote(ctx context.Context, reqBody io.ReadCloser) error {
	decoder := json.NewDecoder(reqBody)
	var body FSNoteCreatedReqBody

	err := decoder.Decode(&body)
	if err != nil {
		log.Println("error decoding request body", err)
		return err
	}

	if containsIssueCmd(body.Data.Text) {
		// Hacky way of grabbing the "sessionID" from the sessionUrl. I'm not the biggest fan of this approach of splitting a string that I don't control.
		sessionId := strings.Split(body.Data.SessionUrl, "/")[6]

		// Check for an existing issue and comment on it.
		potentialExistingIssue, err := inquireExistingIssue(ctx, sessionId)
		if err != nil {
			log.Println("error in inquireExistingIssue", err)
			return err
		}
		if potentialExistingIssue != nil {
			if err := commentOnExistingIssue(ctx, potentialExistingIssue, body.Data.SessionUrl, body.Data.PageInfo.PageUrl, body.Data.Text, body.Data.Author); err != nil {
				log.Println("error creating comment on existing github issue", err)
				return err
			}
		} else {
			// Create the new github issue.
			if err := createGithubIssue(ctx, fmt.Sprintf("Error in session %s", sessionId), body.Data.SessionUrl, body.Data.PageInfo.PageUrl, body.Data.Text, body.Data.Author); err != nil {
				log.Println("error creating github issue", err)
				return err
			}
		}
	}
	return nil
}

func envMust(envVar string) (string, error) {
	if v := os.Getenv(envVar); v == "" {
		return "", errors.New(fmt.Sprintf("error enviorment variable %s cannot be empty", envVar))
	} else {
		return v, nil
	}
}
