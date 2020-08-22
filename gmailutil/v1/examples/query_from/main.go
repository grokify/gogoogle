package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/googleutil/gmailutil/v1"
	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/type/stringsutil"
	omg "github.com/grokify/oauth2more/google"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"google.golang.org/api/gmail/v1"
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
	gs, err := gmailutil.NewGmailService(client)
	if err != nil {
		log.Fatal(err)
	}

	if 1 == 0 {
		labels, err := gmailutil.GetLabelNames(client)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(labels)
	}

	if 1 == 0 {
		rfc822s := []string{
			"list1@example.com",
			"list2@example.com",
			"list3@example.com",
		}
		rfc822sRaw := os.Getenv("EMAIL_ADDRESSES_TO_DELETE")
		if len(rfc822sRaw) > 0 {
			rfc822s = stringsutil.SliceCondenseSpace(strings.Split(rfc822sRaw, ","), true, true)
			fmt.Printf("EMAILS: %s\n", strings.Join(rfc822s, ","))
		}

		deletedCount, gte100Count := gmailutil.DeleteMessagesFrom(gs, rfc822s)

		fmt.Printf("[TOT] DELETED [%v] messages\n", deletedCount)
		fmt.Printf("[TOT] Over 100 [%v] email addresses\n", gte100Count)
	}

	if 1 == 1 {
		msgs, err := GetMessagesByCategory(
			gs, "me", gmailutil.CategoryForums, true)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(msgs)
	}

	fmt.Println("DONE")
}

func GetClient(cfgJson []byte, scopes []string, forceNewToken bool) *http.Client {
	googleClient, err := omg.NewClientFileStoreWithDefaults(
		cfgJson, scopes, forceNewToken)
	if err != nil {
		log.Fatal(errors.Wrap(err, "NewClientFileStoreWithDefaults"))
	}
	return googleClient
}

func GetMessagesByCategory(gs *gmailutil.GmailService, userId, categoryName string, getAll bool) ([]*gmail.Message, error) {
	qOpts := gmailutil.MessagesListQueryOpts{
		Category: categoryName,
	}
	opts := gmailutil.MessagesListOpts{
		Query: qOpts,
	}

	listRes, err := gmailutil.GetMessagesList(gs, opts)
	if err != nil {
		fmt.Println("ERR [%s]", err.Error())
		return []*gmail.Message{}, err
	}
	for _, msg := range listRes.Messages {
		fmtutil.PrintJSON(msg)
		break
	}

	return gmailutil.InflateMessages(gs, userId, listRes.Messages)
}
