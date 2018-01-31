// Go example that covers:
// Quickstart: https://developers.google.com/slides/quickstart/go
// Basic writing: adding a text box to slide: https://developers.google.com/slides/samples/writing
// Using SDK: https://github.com/google/google-api-go-client/blob/master/slides/v1/slides-gen.go
// Creating and Managing Presentations https://developers.google.com/slides/how-tos/presentations
// Adding Shapes and Text to a Slide: https://developers.google.com/slides/how-tos/add-shape#example
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/gotilla/fmt/fmtutil"
	omg "github.com/grokify/oauth2more/google"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"

	"github.com/Iwark/spreadsheet"
)

func loadEnv() error {
	envPaths := []string{}
	envPath := os.Getenv("ENV_PATH")
	if len(envPath) > 0 {
		envPaths = append(envPaths, envPath)
	}
	return godotenv.Load(envPaths...)
}

func main() {
	var forceNewToken bool
	flag.BoolVar(&forceNewToken, "newtoken", false, "Force a new token")
	flag.Parse()

	err := loadEnv()
	if err != nil {
		if strings.Index(err.Error(), "token expired and refresh token is not set") > 0 {
			log.Fatal(fmt.Sprintf("%v - Use option `-newtoken true` to refresh", err.Error()))
		} else {
			log.Fatal(err)
		}
	}

	clientConfig := omg.ClientOauthCliTokenStoreConfig{
		Context:       context.TODO(),
		AppConfig:     []byte(os.Getenv(omg.ClientSecretEnv)),
		Scopes:        []string{sheets.DriveScope, sheets.SpreadsheetsScope},
		TokenFile:     "sheets.googleapis.com-go-quickstart.json",
		ForceNewToken: forceNewToken,
	}

	googleClient, err := omg.NewClientOauthCliTokenStore(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	useIwark := true
	useGoog := false
	if useIwark {
		service := spreadsheet.NewServiceWithClient(googleClient)
		ss, err := service.CreateSpreadsheet(spreadsheet.Spreadsheet{
			Properties: spreadsheet.Properties{
				Title: "spreadsheet title X",
			},
		})

		if err != nil {
			panic(err)
		}

		sheet, err := ss.SheetByIndex(0)
		if err != nil {
			panic(err)
		}

		sheet.Update(3, 2, "Woza2")
		err = sheet.Synchronize()
		if err != nil {
			panic(err)
		}
	}
	if useGoog {
		svc, err := sheets.New(googleClient)
		if err != nil {
			log.Fatal(err)
		}
		sheetsService := sheets.NewSpreadsheetsService(svc)

		ctx := context.Background()
		rb := &sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: "GAPI SHEET",
			},
		}

		resp, err := sheetsService.Create(rb).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)
		}

		fmtutil.PrintJSON(resp)
	}
	fmt.Println("DONE")
}
