package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v35/github"
)

func handleNote(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "handleNote\n")
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
	http.HandleFunc("/handleNote", handleNote)

	http.ListenAndServe(":"+port, nil)

	_ = github.NewClient(nil)
}
