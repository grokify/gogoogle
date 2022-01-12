package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	gu "github.com/grokify/goauth/google"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/urlutil"
	texttospeech "google.golang.org/api/texttospeech/v1beta1"
)

const (
	EnUs    = "en-US"
	Text    = "I like the dreams of the future better than the history of the past."
	Male    = "MALE"
	Female  = "FEMALE"
	Neutral = "NEUTRAL"
	Name    = "en-US-Wavenet-A"
	MP3     = "MP3"
)

func TextSynthesize(ctx context.Context, ttsService *texttospeech.Service) error {
	textService := texttospeech.NewTextService(ttsService)
	synthesizeSpeechRequest := &texttospeech.SynthesizeSpeechRequest{
		AudioConfig: &texttospeech.AudioConfig{
			AudioEncoding: MP3},
		Input: &texttospeech.SynthesisInput{
			Text: Text},
		Voice: &texttospeech.VoiceSelectionParams{
			Name:         Name,
			LanguageCode: EnUs}}

	textSynthesizeCall := textService.Synthesize(synthesizeSpeechRequest)
	textSynthesizeCall.Context(ctx)
	synthesizeSpeechResponse, err := textSynthesizeCall.Do()
	if err != nil {
		return errorsutil.Wrap(err, "TextSynthesize")
	}
	fmtutil.PrintJSON(synthesizeSpeechResponse)

	audio, err := base64.StdEncoding.DecodeString(synthesizeSpeechResponse.AudioContent)
	if err != nil {
		return errorsutil.Wrap(err, "TextSynthesize")
	}
	filename := urlutil.ToSlugLowerString(Text) + "_" + Name + "." + strings.ToLower(MP3)
	err = ioutil.WriteFile(filepath.Join("output", filename), audio, 0644)
	if err != nil {
		return errorsutil.Wrap(err, "TextSynthesize")
	}
	fmt.Printf("WROTE: %v\n", filename)
	return nil
}

func GetVoicesList(ctx context.Context, ttsService *texttospeech.Service) error {
	voiceService := texttospeech.NewVoicesService(ttsService)
	voicesListCall := voiceService.List()
	voicesListCall.LanguageCode(EnUs)
	voicesListCall.Context(ctx)
	listVoicesResponse, err := voicesListCall.Do()
	if err != nil {
		return errorsutil.Wrap(err, "GetVoicesList")
	}
	fmtutil.PrintJSON(listVoicesResponse)
	return nil
}

func main() {
	if err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env"); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	httpClient, err := gu.NewClientFromJWTJSON(
		ctx,
		[]byte(os.Getenv("GOOGLE_SERVICE_ACCOUNT_JWT")),
		texttospeech.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	ttsService, err := texttospeech.New(httpClient)
	if err != nil {
		log.Fatal(err)
	}

	if err = GetVoicesList(ctx, ttsService); err != nil {
		log.Fatal(err)
	}
	if err = TextSynthesize(ctx, ttsService); err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
