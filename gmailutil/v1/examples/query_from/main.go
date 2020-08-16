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
	//Products             string `short:"p" long:"productSlugs" description:"Aha Product Slugs" required:"true"`
	//ReleaseQuarterBegin  int32  `short:"b" long:"begin" description:"Begin Quarter" required:"true"`
	//ReleaseQuarterFinish int32  `short:"f" long:"finish" description:"Finish Quarter" required:"true"`
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

	msgs, err := GetMessagesFrom(client, "listings@redfin.com")
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(msgs.Messages)
	fmt.Printf("NUM_MSGS [%v]\n", len(msgs.Messages))

	fmt.Println("DONE")
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
