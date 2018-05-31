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

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	uu "github.com/grokify/gotilla/net/urlutil"
	gu "github.com/grokify/oauth2more/google"
	"google.golang.org/api/texttospeech/v1beta1"
)

const (
	EnUs    = "en-US"
	Text    = "I like the dreams of the future better than the history of the past."
	Text1   = "Our greatest glory is not in never falling, but in rising every time we fall."
	Male    = "MALE"
	Female  = "FEMALE"
	Neutral = "NEUTRAL"
	Name    = "en-US-Wavenet-E"
	MP3     = "MP3"
)

func TextSynthesize(ctx context.Context, ttsService *texttospeech.Service) {
	textService := texttospeech.NewTextService(ttsService)
	synthesizeSpeechRequest := &texttospeech.SynthesizeSpeechRequest{
		AudioConfig: &texttospeech.AudioConfig{
			AudioEncoding: MP3},
		Input: &texttospeech.SynthesisInput{
			Text: Text},
		Voice: &texttospeech.VoiceSelectionParams{
			Name:         Name,
			LanguageCode: EnUs},
	}
	textSynthesizeCall := textService.Synthesize(synthesizeSpeechRequest)
	textSynthesizeCall.Context(ctx)
	synthesizeSpeechResponse, err := textSynthesizeCall.Do()
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(synthesizeSpeechResponse)
	filename := uu.ToSlugLowerString(Text) + "_" + Name + "." + strings.ToLower(MP3)
	audio, err := base64.StdEncoding.DecodeString(synthesizeSpeechResponse.AudioContent)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join("output", filename), audio, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE: %v\n", filename)
}

func GetVoicesList(ctx context.Context, ttsService *texttospeech.Service) {
	voiceService := texttospeech.NewVoicesService(ttsService)
	voicesListCall := voiceService.List()
	voicesListCall.LanguageCode(EnUs)
	voicesListCall.Context(ctx)
	listVoicesResponse, err := voicesListCall.Do()
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(listVoicesResponse)
}

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	googleJwt := os.Getenv("GOOGLE_SERVICE_ACCOUNT_JWT")
	fmt.Println(googleJwt)

	ctx := context.Background()
	httpClient, err := gu.NewClientFromJWTJSON(
		ctx,
		[]byte(googleJwt),
		texttospeech.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	ttsService, err := texttospeech.New(httpClient)
	if err != nil {
		log.Fatal(err)
	}

	GetVoicesList(ctx, ttsService)
	TextSynthesize(ctx, ttsService)

	fmt.Println("DONE")
}
