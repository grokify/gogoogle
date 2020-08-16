package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grokify/googleutil/gmailutil/v1"
	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	omg "github.com/grokify/oauth2more/google"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

type Options struct {
	EnvFile     string `short:"e" long:"env" description:"Env filepath"`
	NewTokenRaw []bool `short:"n" long:"newtoken" description:"Retrieve new token"`
}

func (opt *Options) NewToken() bool {
	if len(opt.NewTokenRaw) > 0 {
		return true
	}
	return false
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	err = config.LoadDotEnvFirst(opts.EnvFile, os.Getenv("ENV_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	query := gmailutil.MessagesListQueryOpts{From: "foo@example.com"}

	fmt.Printf("%v\n", query)

	client, err := omg.NewClientFileStoreWithDefaults(
		[]byte(os.Getenv(omg.EnvGoogleAppCredentials)),
		[]string{},
		opts.NewToken())
	if err != nil {
		log.Fatal(err)
	}

	labels, err := gmailutil.GetLabelNames(client)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(labels)

	rfc822s := []string{
		"list1@example.com",
		"list2@example.com",
		"list3@example.com",
	}

	for _, rfc822 := range rfc822s {
		ids, err := DeleteMessagesFrom(client, rfc822)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("DELETED [from:%v] [%v] messages\n", rfc822, len(ids))
	}

	fmt.Println("DONE")
}

func DeleteMessagesFrom(client *http.Client, rfc822 string) ([]string, error) {
	ids := []string{}
	msgs, err := GetMessagesFrom(client, rfc822)
	if err != nil {
		return ids, err
	}

	for _, msg := range msgs.Messages {
		ids = append(ids, msg.Id)
	}

	if len(ids) == 0 {
		return ids, nil
	}
	return ids, gmailutil.BatchDeleteMessages(client, []googleapi.CallOption{}, "", ids)
}

func GetMessagesFrom(client *http.Client, rfc822 string) (*gmail.ListMessagesResponse, error) {
	opts := gmailutil.MessagesListOpts{
		Query: gmailutil.MessagesListQueryOpts{
			From: rfc822},
	}

	msgs, err := gmailutil.GetMessagesList(client, []googleapi.CallOption{}, opts)

	return msgs, err
}

func GetClient(cfgJson []byte, scopes []string, forceNewToken bool) *http.Client {
	googleClient, err := omg.NewClientFileStoreWithDefaults(
		cfgJson, scopes, forceNewToken)
	if err != nil {
		log.Fatal(errors.Wrap(err, "NewClientFileStoreWithDefaults"))
	}
	return googleClient
}
