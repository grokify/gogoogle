package gmailutil

import (
	"errors"
	"strings"

	"google.golang.org/api/gmail/v1"
)

func GetMessage(gs *GmailService, userId, messageId string) (*gmail.Message, error) {
	if gs == nil {
		return nil, errors.New("E_NIL_GMAIL_SERVICE")
	}
	return gs.UsersService.Messages.Get(
		strings.TrimSpace(userId),
		strings.TrimSpace(messageId)).
		Do(gs.APICallOptions...)
}
