package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grokify/goauth/google"
	"github.com/grokify/gogoogle/sheetsutil/v4/sheetsmap"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"google.golang.org/api/sheets/v4"
)

func GetSheetsMap() (*sheetsmap.SheetsMap, error) {
	_, err := config.LoadDotEnv([]string{os.Getenv("ENV_PATH"), "./.env"}, -1)
	if err != nil {
		return nil, err
	}

	jwt := os.Getenv("GOOGLE_SERVICE_ACCOUNT_JWT")
	if len(jwt) < 1 {
		return nil, fmt.Errorf("no Google JWT")
	}
	fmt.Println(jwt)

	googleClient, err := google.NewClientFromJWTJSON(
		context.TODO(),
		[]byte(jwt),
		sheets.DriveScope,
		sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	spreadsheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")

	sm, err := sheetsmap.NewSheetsMapIndex(
		googleClient,
		spreadsheetID,
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

	fmtutil.MustPrintJSON(smap.ColumnMapKeyLc)
	fmtutil.MustPrintJSON(smap.ItemMap)

	item, _ := smap.GetOrCreateItem("me@example.com")
	fmtutil.MustPrintJSON(item)

	str, err := smap.UpdateItem(item, "tshirt size", "M", true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	fmt.Println("DONE")
}
