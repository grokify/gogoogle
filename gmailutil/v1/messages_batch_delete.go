package gmailutil

import (
	"errors"
	"fmt"
	"strings"

	gmail "google.golang.org/api/gmail/v1"
)

func BatchDeleteMessages(gs *GmailService, userID string, messageIDs []string) error {
	if gs == nil {
		return errors.New("nil `GmailService`")
	}

	userID = strings.TrimSpace(userID)
	if len(userID) == 0 {
		userID = "me"
	}

	return gs.UsersService.Messages.BatchDelete(
		userID,
		&gmail.BatchDeleteMessagesRequest{Ids: messageIDs}).
		Do(gs.APICallOptions...)
}

func DeleteMessagesFrom(gs *GmailService, rfc822s []string) (int, int, error) {
	deletedCount := 0
	gte100Count := 0
	rfc822Count := len(rfc822s)
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
