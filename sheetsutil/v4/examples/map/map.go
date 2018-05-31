package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/googleutil/sheetsutil/v4/sheetsmap"
	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	gu "github.com/grokify/oauth2more/google"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

func GetSheetsMap() (*sheetsmap.SheetsMap, error) {
	log.SetLevel(log.DebugLevel)

	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		return nil, err
	}

	jwt := os.Getenv("GOOGLE_SERVICE_ACCOUNT_JWT")
	if len(jwt) < 1 {
		return nil, fmt.Errorf("No Google JWT")
	}
	fmt.Println(jwt)

	googleClient, err := gu.NewClientFromJWTJSON(
		context.TODO(),
		[]byte(jwt),
		sheets.DriveScope,
		sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	spreadsheetId := os.Getenv("GOOGLE_SPREADSHEET_ID")

	sm, err := sheetsmap.NewSheetsMapIndex(
		googleClient,
		spreadsheetId,
		uint(0))
	return &sm, err
}

func main() {
	smap, err := GetSheetsMap()
	if err != nil {
		log.Fatal(err)
	}

	err = smap.ReadColumns()
	if err != nil {
		log.Fatal(err)
	}

	err = smap.FullRead()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("VAL %v\n", smap.Sheet.Rows[1][0].Value)

	fmtutil.PrintJSON(smap.ColumnMapKeyLc)
	fmtutil.PrintJSON(smap.ItemMap)

	item, _ := smap.GetOrCreateItem("me@example.com")
	fmtutil.PrintJSON(item)

	str, err := smap.UpdateItem(item, "tshirt size", "M", true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	fmt.Println("DONE")
}
