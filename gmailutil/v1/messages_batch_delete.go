package gmailutil

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/gmail/v1"
)

func BatchDeleteMessages(gs *GmailService, userId string, messageIds []string) error {
	userId = strings.TrimSpace(userId)
	if len(userId) == 0 {
		userId = "me"
	}

	usersMessagesBatchDeleteCall := gs.UsersService.Messages.BatchDelete(
		userId,
		&gmail.BatchDeleteMessagesRequest{Ids: messageIds})

	return usersMessagesBatchDeleteCall.Do(gs.APICallOptions...)
}

func DeleteMessagesFrom(gs *GmailService, rfc822s []string) (int, int) {
	deletedCount := 0
	gte100Count := 0
	for i, rfc822 := range rfc822s {
		ids, err := deleteMessagesFromSingle(gs, rfc822)
		if err != nil {
			log.Fatal(err)
		}
		numDeleted := len(ids)
		alert := ""
		if numDeleted >= 100 {
			alert = " (>100)"
			gte100Count++
		}
		fmt.Printf("[%d] DELETED [%v]%s messages [from:%v]\n", i+1, numDeleted, alert, rfc822)
		deletedCount += numDeleted
	}
	return deletedCount, gte100Count
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
