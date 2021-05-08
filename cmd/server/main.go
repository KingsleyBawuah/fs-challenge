package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/google/go-github/v35/github"
)

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

// containsIssueCmd determines if the note text contains the string #issue which indicates this note should create an issue against the site repo.
func containsIssueCmd(text string) bool {
	matchIssueCmd := `\W(\#(issue)+\b)`

	match, err := regexp.MatchString(matchIssueCmd, text)
	if err != nil {
		log.Panicln("error determining if note text contains issue command", err)
	}

	return match
}

func handleNote(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		decoder := json.NewDecoder(req.Body)
		var body FSNoteCreatedReqBody
		err := decoder.Decode(&body)
		if err != nil {
			log.Panicln("error decoding request body", err)
		}

		if containsIssueCmd(body.Data.Text) {
			// Create the github issue.
			fmt.Fprintf(w, "Yo that's a cmd")
		} else {
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

func main() {
	port, err := envMust("PORT")
	if err != nil {
		os.Exit(1)
	}

	// Set up HTTP Request Handlers.
	// TODO: Handle errors here better and defer closing the connection if needed.
	http.HandleFunc("/handleNote", handleNote)

	http.ListenAndServe(":"+port, nil)

	log.Printf("Server started at port %s", port)

	_ = github.NewClient(nil)
}
