package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
	"google.golang.org/api/option"

	//gu "github.com/grokify/oauth2more/google"

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

func NewRecognitionAudio(data []byte) *speechpb.RecognitionAudio {
	return &speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
	}
}

func NewRecognitionAudioFile(file string) (*speechpb.RecognitionAudio, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return NewRecognitionAudio(data), nil
}

func RecognizeResponseTextFirst(resp *speechpb.RecognizeResponse, threshold float32) (string, error) {
	//highestThreshold := 0.0
	//highestThresholdTranscript := ""
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			trimmed := strings.TrimSpace(alt.Transcript)
			if len(trimmed) > 0 {
				if threshold == 0 ||
					(threshold > 0 && alt.Confidence >= threshold) {
					return trimmed, nil
				}
			}
		}
	}
	return "", fmt.Errorf("No responses found")
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
	audio, err := NewRecognitionAudioFile(args.AudioFile)
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
	firstText, err := RecognizeResponseTextFirst(resp, 0.5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FIRST: [%v]\n", firstText)
	fmt.Println("S5")
	fmt.Println("DONE")
}
