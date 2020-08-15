package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grokify/googleutil/gmailutil/v1"
	"github.com/grokify/gotilla/config"
	omg "github.com/grokify/oauth2more/google"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"google.golang.org/api/gmail/v1"
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

	query := gmailutil.MessageListQueryOpts{From: "foo@example.com"}

	fmt.Printf("%v\n", query)

	client := GetClient(opts)

	Quickstart(client)

	fmt.Println("DONE")
}

func GetClient(opts Options) *http.Client {
	googleClient, err := omg.NewClientFileStoreWithDefaults(
		[]byte(os.Getenv(omg.EnvGoogleAppCredentials)),
		[]string{gmail.GmailReadonlyScope},
		opts.NewToken())
	if err != nil {
		log.Fatal(errors.Wrap(err, "NewClientFileStoreWithDefaults"))
	}
	return googleClient
}

func Quickstart(client *http.Client) {
	// https://developers.google.com/gmail/api/quickstart/go
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}
