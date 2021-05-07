package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v35/github"
	"github.com/joho/godotenv"
)

func handleNote(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "handleNote\n")
}

func main() {
	err := godotenv.Load()

	if err != nil {
		//TODO: DO not panic.
		log.Panic(err)
	}
	port := os.Getenv("PORT")
	fmt.Println("PORT is", port)

	http.HandleFunc("/handleNote", handleNote)

	//TODO: Figure out how to read env for port number.
	http.ListenAndServe(":"+port, nil)

	_ = github.NewClient(nil)
}
