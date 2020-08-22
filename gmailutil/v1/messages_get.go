package gmailutil

import (
	"strings"

	"google.golang.org/api/gmail/v1"
)

func GetMessage(gs *GmailService, userId, messageId string) (*gmail.Message, error) {
	userId = strings.TrimSpace(userId)
	messageId = strings.TrimSpace(messageId)

	userMessagesListCall := gs.UsersService.Messages.Get(
		userId, messageId)
	return userMessagesListCall.Do(gs.APICallOptions...)
}
