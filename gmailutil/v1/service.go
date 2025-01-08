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
	if service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client)); err != nil {
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

// Send is a helpfer for https://pkg.go.dev/google.golang.org/api/gmail/v1#UsersMessagesService.Send
func (gs GmailService) Send(ctx context.Context, from string, msg mailutil.MessageWriter, opts ...googleapi.CallOption) (*gmail.Message, error) {
	gmsg := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(msg.String()))}
	call := gs.UsersService.Messages.Send(from, gmsg)
	call = call.Context(ctx)
	return call.Do(opts...)
}
