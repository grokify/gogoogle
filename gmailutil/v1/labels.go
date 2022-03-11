package gmailutil

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func GetLabelNames(client *http.Client) ([]string, error) {
	// https://developers.google.com/gmail/api/quickstart/go
	labels := []string{}
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return labels, fmt.Errorf("unable to retrieve gmail client: err [%v]", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		return labels, fmt.Errorf("unable to retrieve labels: err [%v]", err)
	}
	for _, l := range r.Labels {
		labels = append(labels, l.Name)
	}
	return labels, nil
}
