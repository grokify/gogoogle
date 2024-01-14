package gmailutil

import (
	"fmt"
	"strings"

	gmail "google.golang.org/api/gmail/v1"
)

func (mapi *MessagesAPI) BatchDeleteMessages(userID string, messageIDs []string) error {
	if mapi.GmailService == nil {
		return ErrGmailServiceCannotBeNil
	}

	userID = strings.TrimSpace(userID)
	if len(userID) == 0 {
		return ErrGmailUserIDCannotBeEmpty
	}

	return mapi.GmailService.UsersService.Messages.BatchDelete(
		userID,
		&gmail.BatchDeleteMessagesRequest{Ids: messageIDs}).
		Do(mapi.GmailService.APICallOptions...)
}

func (mapi *MessagesAPI) DeleteMessagesFrom(rfc822s []string) (int, int, error) {
	if mapi.GmailService == nil {
		return -1, -1, ErrGmailServiceCannotBeNil
	}

	deletedCount := 0
	gte100Count := 0
	rfc822Count := len(rfc822s)
	for i, rfc822 := range rfc822s {
		ids, err := mapi.deleteMessagesFromSingle(rfc822)
		if err != nil {
			return deletedCount, gte100Count, err
		}
		numDeleted := len(ids)
		alert := ""
		if numDeleted >= 100 {
			alert = " (>100)"
			gte100Count++
		}
		for numDeleted >= 100 {
			ids, err := mapi.deleteMessagesFromSingle(rfc822)
			if err != nil {
				return deletedCount, gte100Count, err
			}
			numDeleted = len(ids)
		}
		fmt.Printf("[%d/%d] DELETED [%v]%s messages [from:%v]\n", i+1, rfc822Count, numDeleted, alert, rfc822)
		fmt.Printf("{\"addressNum\":%d, \"addressTotal\": %d, \"deletedCount\": %d}", i+1, rfc822Count, numDeleted)
		/*log.Info().
		Int("address_num", i+1).
		Int("address_total", rfc822Count).
		Int("deleted_count", numDeleted)*/
		deletedCount += numDeleted
	}
	return deletedCount, gte100Count, nil
}

func (mapi *MessagesAPI) deleteMessagesFromSingle(rfc822 string) ([]string, error) {
	ids := []string{}
	listRes, err := mapi.GetMessagesFrom(rfc822)
	if err != nil {
		return ids, err
	}

	for _, msg := range listRes.Messages {
		ids = append(ids, msg.Id)
	}

	if len(ids) == 0 {
		return ids, nil
	}
	return ids, mapi.BatchDeleteMessages(UserIDMe, ids)
}
