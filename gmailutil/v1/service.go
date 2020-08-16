package gmailutil

import (
	"net/http"

	"google.golang.org/api/gmail/v1"
)

func NewUsersService(client *http.Client) (*gmail.UsersService, error) {
	service, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	return gmail.NewUsersService(service), nil
}
