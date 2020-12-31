package gmailutil

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/gmail/v1"
)

func BatchDeleteMessages(gs *GmailService, userId string, messageIds []string) error {
	if gs == nil {
		return errors.New("E_NIL_GMAIL_SERVICE")
	}

	userId = strings.TrimSpace(userId)
	if len(userId) == 0 {
		userId = "me"
	}

	return gs.UsersService.Messages.BatchDelete(
		userId,
		&gmail.BatchDeleteMessagesRequest{Ids: messageIds}).
		Do(gs.APICallOptions...)
}

func DeleteMessagesFrom(gs *GmailService, rfc822s []string) (int, int, error) {
	deletedCount := 0
	gte100Count := 0
	for i, rfc822 := range rfc822s {
		ids, err := deleteMessagesFromSingle(gs, rfc822)
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
			ids, err := deleteMessagesFromSingle(gs, rfc822)
			if err != nil {
				return deletedCount, gte100Count, err
			}
			numDeleted = len(ids)
		}
		fmt.Printf("[%d] DELETED [%v]%s messages [from:%v]\n", i+1, numDeleted, alert, rfc822)
		deletedCount += numDeleted
	}
	return deletedCount, gte100Count, nil
}

func deleteMessagesFromSingle(gs *GmailService, rfc822 string) ([]string, error) {
	ids := []string{}
	listRes, err := GetMessagesFrom(gs, rfc822)
	if err != nil {
		return ids, err
	}

	for _, msg := range listRes.Messages {
		ids = append(ids, msg.Id)
	}

	if len(ids) == 0 {
		return ids, nil
	}
	return ids, BatchDeleteMessages(gs, "", ids)
}
