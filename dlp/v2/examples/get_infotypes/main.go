package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grokify/mogo/fmt/fmtutil"

	"github.com/jessevdk/go-flags"
	"google.golang.org/api/option"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

/*

https://godoc.org/cloud.google.com/go/dlp/apiv2
https://cloud.google.com/dlp/docs/auth
https://dlp.googleapis.com/v2/infoTypes

Note: Granting the user the Owner role is not necessary.

This is based on the following Quickstart

https://github.com/GoogleCloudPlatform/golang-samples/blob/master/dlp/dlp_quickstart/quickstart.go

*/

type Args struct {
	// Service Account Credentials File
	CredentialsFile string `short:"c" long:"credentials" description:"Path to crdentials file." required:"true"`
}

func main() {
	args := Args{}

	_, err := flags.Parse(&args)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	opts := option.WithCredentialsFile(args.CredentialsFile)

	client, err := dlp.NewClient(ctx, opts)
	if err != nil {
		log.Fatalf("error creating DLP client: %v", err)
	}

	res, err := client.ListInfoTypes(ctx, &dlppb.ListInfoTypesRequest{})
	if err != nil {
		log.Fatalf("error retrieving InfoTypes: %v", err)
	}

	fmtutil.PrintJSON(res)

	for i, t := range res.InfoTypes {
		fmt.Printf("%v) %v\n", i+1, t.Name)
	}

	fmt.Println("DONE")
}
