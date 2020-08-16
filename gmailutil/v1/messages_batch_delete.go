package gmailutil

import (
	"net/http"
	"strings"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

func BatchDeleteMessages(client *http.Client, apiOpts []googleapi.CallOption, userId string, messageIds []string) error {
	usersService, err := NewUsersService(client)
	if err != nil {
		return err
	}
	userId = strings.TrimSpace(userId)
	if len(userId) == 0 {
		userId = "me"
	}

	usersMessagesBatchDeleteCall := usersService.Messages.BatchDelete(
		userId,
		&gmail.BatchDeleteMessagesRequest{Ids: messageIds})

	return usersMessagesBatchDeleteCall.Do(apiOpts...)
}
