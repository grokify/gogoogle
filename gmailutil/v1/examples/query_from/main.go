package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	omg "github.com/grokify/goauth/google"
	gmailutil "github.com/grokify/gogoogle/gmailutil/v1"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	EnvFile     string `short:"e" long:"env" description:"Env filepath"`
	NewTokenRaw []bool `short:"n" long:"newtoken" description:"Retrieve new token"`
}

func (opt *Options) NewToken() bool {
	return len(opt.NewTokenRaw) > 0
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = config.LoadDotEnv([]string{opts.EnvFile, os.Getenv("ENV_PATH")}, 1)
	if err != nil {
		log.Fatal(err)
	}

	query := gmailutil.MessagesListQueryOpts{From: "foo@example.com"}

	fmt.Printf("%v\n", query)

	client, err := omg.NewClientFileStoreWithDefaults(
		context.Background(),
		[]byte(os.Getenv(omg.EnvGoogleAppCredentials)),
		[]string{},
		opts.NewToken())
	if err != nil {
		log.Fatal(err)
	}
	gs, err := gmailutil.NewGmailService(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}

	if 1 == 0 {
		labels, err := gmailutil.GetLabelNames(client)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.MustPrintJSON(labels)
	}

	if 1 == 1 {
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

		deletedCount, gte100Count, err := gs.MessagesAPI.DeleteMessagesFrom(rfc822s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[TOT] DELETED [%v] messages\n", deletedCount)
		fmt.Printf("[TOT] Over 100 [%v] email addresses\n", gte100Count)
	}

	if 1 == 0 {
		msgs, err := gs.MessagesAPI.GetMessagesByCategory(
			gmailutil.UserIDMe, gmailutil.CategoryForums, true)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.MustPrintJSON(msgs)
	}

	fmt.Println("DONE")
}

func GetClient(cfgJSON []byte, scopes []string, forceNewToken bool) *http.Client {
	googleClient, err := omg.NewClientFileStoreWithDefaults(
		context.Background(), cfgJSON, scopes, forceNewToken)
	if err != nil {
		log.Fatal(errorsutil.Wrap(err, "NewClientFileStoreWithDefaults"))
	}
	return googleClient
}

/*
func GetMessagesByCategory(gs *gmailutil.GmailService, userId, categoryName string, getAll bool) ([]*gmail.Message, error) {
	qOpts := gmailutil.MessagesListQueryOpts{
		Category: categoryName,
	}
	opts := gmailutil.MessagesListOpts{
		Query: qOpts,
	}

	listRes, err := gmailutil.GetMessagesList(gs, opts)
	if err != nil {
		fmt.Printf("ERR [%s]", err.Error())
		return []*gmail.Message{}, err
	}
	for _, msg := range listRes.Messages {
		fmtutil.PrintJSON(msg)
		break
	}

	return gmailutil.InflateMessages(gs, userId, listRes.Messages)
}
*/
