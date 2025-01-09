package iwark

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Iwark/spreadsheet"
	"github.com/grokify/goauth"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gogoogle/docsutil"
	"github.com/grokify/mogo/type/stringsutil"
)

var ErrSheetIDRequired = errors.New("sheet id is required")

func ParseTableFromSheetIDCredentials(ctx context.Context, creds *goauth.CredentialsSet, credsKey string, sheetID string, sheetIdx, headerRows uint) (*table.Table, error) {
	if c, err := creds.Get(credsKey); err != nil {
		return nil, err
	} else if clt, err := c.NewClient(ctx); err != nil {
		return nil, err
	} else {
		return ParseTableFromSheetIDClient(clt, sheetID, 0, 1)
	}
}

func ParseTableFromSheetIDClient(client *http.Client, sheetID string, sheetIdx, headerRows uint64) (*table.Table, error) {
	if strings.Contains(sheetID, "/") {
		id, _, err := docsutil.ParseDocsURL(sheetID, docsutil.DocSlugSpreadsheet)
		if err == nil && id != "" {
			sheetID = id
		}
	}
	sheetID = strings.TrimSpace(sheetID)
	if sheetID == "" {
		return nil, ErrSheetIDRequired
	}
	service := spreadsheet.NewServiceWithClient(client)
	if ss, err := service.FetchSpreadsheet(sheetID); err != nil {
		return nil, err
	} else {
		return ParseTableFromSpreadsheet(ss, sheetIdx, headerRows)
	}
}

func ParseTableFromSpreadsheet(ss spreadsheet.Spreadsheet, sheetIdx, headerRows uint64) (*table.Table, error) {
	if s, err := ss.SheetByIndex(uint(sheetIdx)); err != nil {
		return nil, err
	} else {
		return ParseTableFromSheet(s, headerRows), nil
	}
}

func ParseTableFromSheet(s *spreadsheet.Sheet, headerRows uint64) *table.Table {
	cols, rows := ParseDataFromSheet(s, headerRows)
	tbl := table.NewTable("")
	tbl.Columns = cols
	tbl.Rows = rows
	return &tbl
}

func ParseDataFromSheet(s *spreadsheet.Sheet, headerRows uint64) ([]string, [][]string) {
	var cols []string
	var rows [][]string
	for i, srow := range s.Rows {
		var row []string
		for _, scell := range srow {
			row = append(row, scell.Value)
		}
		if headerRows > 0 && uint64(i) < headerRows {
			if i == 0 {
				cols = row
			}
		} else {
			rows = append(rows, row)
		}
	}
	cols = stringsutil.SliceTrimSpace(cols, false)
	return cols, rows
}
