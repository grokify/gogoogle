package gmailutil

import (
	"fmt"
	"strings"

	gmail "google.golang.org/api/gmail/v1"
)

func (mapi *MessagesAPI) GetMessage(userID, messageID string) (*gmail.Message, error) {
	if mapi.GmailService == nil {
		return nil, ErrGmailServiceCannotBeNil
	}
	return mapi.GmailService.UsersService.Messages.Get(
		strings.TrimSpace(userID),
		strings.TrimSpace(messageID)).
		Do(mapi.GmailService.APICallOptions...)
}

func (mapi *MessagesAPI) GetMessagesByCategory(userID, categoryName string, getAll bool) ([]*gmail.Message, error) {
	if mapi.GmailService == nil {
		return nil, ErrGmailServiceCannotBeNil
	}
	qOpts := MessagesListQueryOpts{
		Category: categoryName,
	}
	opts := MessagesListOpts{
		Query: qOpts,
	}

	listRes, err := mapi.GetMessagesList(opts)
	if err != nil {
		fmt.Printf("ERR [%s]", err.Error())
		return []*gmail.Message{}, err
	}

	return mapi.InflateMessages(userID, listRes.Messages)
}
