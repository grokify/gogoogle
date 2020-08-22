package gmailutil

import (
	"net/http"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

func NewUsersService(client *http.Client) (*gmail.UsersService, error) {
	service, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	return gmail.NewUsersService(service), nil
}

type GmailService struct {
	httpClient     *http.Client
	Service        *gmail.Service
	UsersService   *gmail.UsersService
	APICallOptions []googleapi.CallOption
}

func NewGmailService(client *http.Client) (*GmailService, error) {
	gs := &GmailService{
		httpClient:     client,
		APICallOptions: []googleapi.CallOption{}}
	service, err := gmail.New(client)
	if err != nil {
		return gs, err
	}
	gs.Service = service
	gs.UsersService = gmail.NewUsersService(service)
	return gs, nil
}
