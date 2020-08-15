package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/google-api-go-client/gmail/v1"
	"github.com/grokify/googleutil/gmailutil/v1"
)

func main() {
	query := gmailutil.MessageListQueryOpts{From: "foo@example.com"}

	client := http.Client{}

	Quickstart(client)

	fmt.Println("DONE")
}

func Quickstart(client *http.Client) {
	// https://developers.google.com/gmail/api/quickstart/go
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}
