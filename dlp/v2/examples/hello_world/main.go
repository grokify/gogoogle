package main

import (
	"context"
	"fmt"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
	"google.golang.org/api/option"

	//dlp "cloud.google.com/go/dlp/v2"
	gu "github.com/grokify/goauth/google"
	dlpu "github.com/grokify/googleutil/dlp/v2"

	//dlp "google.golang.org/api/dlp/v2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

/*

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
		log.Fatal(err)
	}

	credsContainer, err := gu.CredentialsContainerFromFile(args.CredentialsFile)
	if err != nil {
		log.Fatal(err)
	}
	creds := credsContainer.Credentials()

	projectID := creds.ProjectID
	input := "Hello World 680-26-5240"

	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: []*dlppb.InfoType{
			{
				Name: dlpu.InfoTypeAllBasic,
			},
		},

		MinLikelihood: dlppb.Likelihood_POSSIBLE,

		Limits: &dlppb.InspectConfig_FindingLimits{
			MaxFindingsPerRequest: int32(0),
		},
		IncludeQuote: true}

	req := &dlppb.InspectContentRequest{
		Parent:        "projects/" + projectID,
		InspectConfig: inspectConfig,
		Item:          dlpu.NewContentDataItemSimple(input)}

	fmtutil.PrintJSON(req)

	// Run request.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	findings := resp.GetResult().GetFindings()
	if len(findings) == 0 {
		fmt.Println("No findings.")
	}
	fmt.Println("Findings:")
	for _, f := range findings {
		if inspectConfig.IncludeQuote {
			fmt.Println("\tQuote: ", f.GetQuote())
		}
		fmt.Println("\tInfo type: ", f.GetInfoType().GetName())
		fmt.Println("\tLikelihood: ", f.GetLikelihood())
	}

	fmt.Println("DONE")
}
