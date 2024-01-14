package gmailutil

import (
	"context"
	"errors"
	"net/http"

	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const UserIDMe = "me"

var (
	ErrGmailServiceCannotBeNil  = errors.New("gmail service cannot be nil")
	ErrGmailUserIDCannotBeEmpty = errors.New("gmail userid cannot be empty")
)

func NewUsersService(client *http.Client) (*gmail.UsersService, error) {
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return gmail.NewUsersService(service), nil
}

type GmailService struct {
	httpClient     *http.Client
	APICallOptions []googleapi.CallOption
	Service        *gmail.Service
	UsersService   *gmail.UsersService
	MessagesAPI    MessagesAPI
}

func NewGmailService(client *http.Client) (*GmailService, error) {
	gs := &GmailService{
		httpClient:     client,
		APICallOptions: []googleapi.CallOption{}}
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return gs, err
	}
	gs.Service = service
	gs.UsersService = gmail.NewUsersService(service)
	gs.MessagesAPI = MessagesAPI{GmailService: gs}
	return gs, nil
}

type MessagesAPI struct {
	GmailService *GmailService
}
