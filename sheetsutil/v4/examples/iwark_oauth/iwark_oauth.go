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

	omg "github.com/grokify/goauth/google"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"

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
			log.Fatalf("%v - Use option `-newtoken true` to refresh", err.Error())
		} else {
			log.Fatal(err)
		}
	}

	clientConfig := omg.ClientOAuthCLITokenStoreConfig{
		Context:       context.Background(),
		AppConfig:     []byte(os.Getenv(omg.ClientSecretEnv)),
		Scopes:        []string{sheets.DriveScope, sheets.SpreadsheetsScope},
		TokenFile:     "sheets.googleapis.com-go-quickstart.json",
		ForceNewToken: forceNewToken,
	}

	hclient, err := omg.NewClientOAuthCLITokenStore(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	useIwark := true
	useGoog := false
	if useIwark {
		service := spreadsheet.NewServiceWithClient(hclient)
		ss, err := service.CreateSpreadsheet(spreadsheet.Spreadsheet{
			Properties: spreadsheet.Properties{
				Title: "spreadsheet title X",
			},
		})

		if err != nil {
			log.Fatal(err)
		}

		sheet, err := ss.SheetByIndex(0)
		if err != nil {
			panic(err)
		}

		err = service.ExpandSheet(sheet, 20, 10)
		if err != nil {
			log.Fatal(err)
		}

		sheet.Update(3, 2, "Woza2")
		err = sheet.Synchronize()
		if err != nil {
			log.Fatal(err)
		}
	}

	if useGoog {
		svc, err := sheets.NewService(context.Background(), option.WithHTTPClient(hclient))
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

		fmtutil.MustPrintJSON(resp)
	}
	fmt.Println("DONE")
}
