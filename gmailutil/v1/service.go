package gmailutil

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/grokify/mogo/net/mailutil"
	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

var (
	ErrGmailServiceCannotBeNil              = errors.New("gmail service cannot be nil")
	ErrGmailUsersServiceCannotBeNil         = errors.New("gmail users service cannot be nil")
	ErrGmailUsersServiceMessagesCannotBeNil = errors.New("gmail users service messages cannot be nil")
	ErrGmailUserIDCannotBeEmpty             = errors.New("gmail userid cannot be empty")
	ErrHTTPClientCannotBeNil                = errors.New("http.client cannot be nil")
	ErrMessagesAPIGmailServiceCannotBeNil   = errors.New("messages api.gmail service cannot be nil")
)

func NewUsersService(client *http.Client) (*gmail.UsersService, error) {
	if client == nil {
		return nil, ErrHTTPClientCannotBeNil
	} else if service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client)); err != nil {
		return nil, err
	} else {
		return gmail.NewUsersService(service), nil
	}
}

type GmailService struct {
	httpClient     *http.Client
	APICallOptions []googleapi.CallOption
	Service        *gmail.Service
	UsersService   *gmail.UsersService
	MessagesAPI    MessagesAPI
}

func NewGmailService(ctx context.Context, client *http.Client) (*GmailService, error) {
	gs := &GmailService{
		httpClient:     client,
		APICallOptions: []googleapi.CallOption{}}
	if service, err := gmail.NewService(ctx, option.WithHTTPClient(client)); err != nil {
		return nil, err
	} else {
		gs.Service = service
	}
	gs.UsersService = gmail.NewUsersService(gs.Service)
	gs.MessagesAPI = MessagesAPI{GmailService: gs}
	return gs, nil
}

type MessagesAPI struct {
	GmailService *GmailService
}

func (gs GmailService) validateConfig() error {
	if gs.httpClient == nil {
		return ErrHTTPClientCannotBeNil
	} else if gs.Service == nil {
		return ErrGmailServiceCannotBeNil
	} else if gs.UsersService == nil {
		return ErrGmailUsersServiceCannotBeNil
	} else if gs.UsersService.Messages == nil {
		return ErrGmailUsersServiceMessagesCannotBeNil
	} else if gs.MessagesAPI.GmailService == nil {
		return ErrMessagesAPIGmailServiceCannotBeNil
	} else {
		return nil
	}
}

// Send is a helpfer for https://pkg.go.dev/google.golang.org/api/gmail/v1#UsersMessagesService.Send
func (gs GmailService) Send(ctx context.Context, from string, msg mailutil.MessageWriter, opts ...googleapi.CallOption) (*gmail.Message, error) {
	if err := gs.validateConfig(); err != nil {
		return nil, err
	}
	gmsg := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(msg.String()))}
	call := gs.UsersService.Messages.Send(from, gmsg)
	call = call.Context(ctx)
	return call.Do(opts...)
}
