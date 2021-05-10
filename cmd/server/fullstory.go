package main

import (
	"log"
	"regexp"
	"time"
)

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

// containsIssueCmd determines if the note text from fullstory contains the string #issue which indicates this note should create an issue against the site repo.
func containsIssueCmd(text string) bool {
	matchIssueCmd := `\W(\#(issue)+\b)`

	match, err := regexp.MatchString(matchIssueCmd, text)
	if err != nil {
		log.Panicln("error determining if note text contains issue command", err)
	}

	return match
}
