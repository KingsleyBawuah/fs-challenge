package main

import "github.com/google/go-github/v35/github"

func main() {
	_ = github.NewClient(nil)
}
