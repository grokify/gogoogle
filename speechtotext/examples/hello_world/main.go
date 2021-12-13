package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
	"google.golang.org/api/option"

	"github.com/grokify/googleutil/speechtotext"
	//gu "github.com/grokify/goauth/google"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

/*

This is based on the following Quickstart

https://github.com/GoogleCloudPlatform/golang-samples/blob/master/dlp/dlp_quickstart/quickstart.go

*/

type Args struct {
	// Service Account Credentials File
	CredentialsFile string `short:"c" long:"credentials" description:"Path to crdentials file." required:"true"`

	// Service Account Credentials File
	AudioFile string `short:"f" long:"file" description:"Path to audio file." required:"true"`
}

func main() {
	args := Args{}

	_, err := flags.Parse(&args)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	client, err := speech.NewClient(ctx,
		option.WithCredentialsFile(args.CredentialsFile))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("S1")
	// ffmpeg -i input.mp3 output.flac
	audio, err := speechtotext.NewRecognitionAudioFile(args.AudioFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("S2")

	config := &speechpb.RecognitionConfig{
		//Encoding:        speechpb.RecognitionConfig_ENCODING_UNSPECIFIED, //
		//Encoding:        speechpb.RecognitionConfig_LINEAR16,
		//SampleRateHertz: 16000,
		LanguageCode: "en-US"}

	req := &speechpb.RecognizeRequest{
		Config: config,
		Audio:  audio}

	fmt.Println("S3")
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("S4")
	fmtutil.PrintJSON(resp)

	// Prints the results.
	for i, result := range resp.Results {
		fmt.Printf("RES [%v]\n", i)
		fmtutil.PrintJSON(result)
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
	firstText, err := speechtotext.RecognizeResponseTextFirst(resp, 0.5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FIRST: [%v]\n", firstText)
	fmt.Println("S5")
	fmt.Println("DONE")
}
