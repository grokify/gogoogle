package gmailutil

import (
	"fmt"
	"net/http"

	"google.golang.org/api/gmail/v1"
)

func GetLabelNames(client *http.Client) ([]string, error) {
	// https://developers.google.com/gmail/api/quickstart/go
	labels := []string{}
	srv, err := gmail.New(client)
	if err != nil {
		return labels, fmt.Errorf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		return labels, fmt.Errorf("Unable to retrieve labels: %v", err)
	}
	for _, l := range r.Labels {
		labels = append(labels, l.Name)
	}
	return labels, nil
}
