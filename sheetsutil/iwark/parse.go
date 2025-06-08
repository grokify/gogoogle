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

var (
	ErrSheetIDRequired  = errors.New("sheet id is required")
	ErrSheetCannotBeNil = errors.New("sheet cannot be nil")
)

func ReadTableFromCredentialsSetFile(ctx context.Context, credsFile, credsKey string, sheetID string, sheetIdx uint, headerRows uint32) (*table.Table, error) {
	if creds, err := goauth.ReadCredentialsFromSetFile(credsFile, credsKey, true); err != nil {
		return nil, err
	} else {
		return ReadTableFromCredentials(ctx, creds, sheetID, sheetIdx, headerRows)
	}
}

func ReadTableFromCredentialsSet(ctx context.Context, credsSet *goauth.CredentialsSet, credsKey string, sheetID string, sheetIdx uint, headerRows uint32) (*table.Table, error) {
	if creds, err := credsSet.Get(credsKey); err != nil {
		return nil, err
	} else {
		return ReadTableFromCredentials(ctx, creds, sheetID, sheetIdx, headerRows)
	}
}

func ReadTableFromCredentials(ctx context.Context, creds goauth.Credentials, sheetID string, sheetIdx uint, headerRows uint32) (*table.Table, error) {
	if clt, err := creds.NewClient(ctx); err != nil {
		return nil, err
	} else {
		return ReadTableFromClient(clt, sheetID, 0, 1)
	}
}

func ReadTableFromClient(client *http.Client, sheetID string, sheetIdx uint, headerRows uint32) (*table.Table, error) {
	if ss, err := ReadSpreadsheetFromClient(client, sheetID); err != nil {
		return nil, err
	} else {
		return ParseTableFromSpreadsheet(*ss, sheetIdx, headerRows)
	}
}

func ReadSheetFromCredentialsSetFile(ctx context.Context, credsFile, credsKey string, sheetID string, sheetIdx uint) (*spreadsheet.Sheet, error) {
	if creds, err := goauth.ReadCredentialsFromSetFile(credsFile, credsKey, true); err != nil {
		return nil, err
	} else if client, err := creds.NewClient(ctx); err != nil {
		return nil, err
	} else {
		return ReadSheetFromClient(client, sheetID, sheetIdx)
	}
}

func ReadSheetFromClient(client *http.Client, sheetID string, sheetIdx uint) (*spreadsheet.Sheet, error) {
	if ss, err := ReadSpreadsheetFromClient(client, sheetID); err != nil {
		return nil, err
	} else {
		return ss.SheetByIndex(sheetIdx)
	}
}

func ReadSpreadsheetFromClient(client *http.Client, sheetID string) (*spreadsheet.Spreadsheet, error) {
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
		return &ss, nil
	}
}

func ParseTableFromSpreadsheet(ss spreadsheet.Spreadsheet, sheetIdx uint, headerRows uint32) (*table.Table, error) {
	if s, err := ss.SheetByIndex(sheetIdx); err != nil {
		return nil, err
	} else {
		return ParseTableFromSheet(s, headerRows)
	}
}

func ParseTableFromSheet(s *spreadsheet.Sheet, headerRows uint32) (*table.Table, error) {
	if s == nil {
		return nil, ErrSheetCannotBeNil
	}
	cols, rows := ParseDataFromSheet(s, headerRows)
	tbl := table.NewTable("")
	tbl.Columns = cols
	tbl.Rows = rows
	return &tbl, nil
}

func ParseDataFromSheet(s *spreadsheet.Sheet, headerRows uint32) ([]string, [][]string) {
	var cols []string
	var rows [][]string
	for i, srow := range s.Rows {
		var row []string
		for _, scell := range srow {
			row = append(row, scell.Value)
		}
		if headerRows > 0 && i < int(headerRows) {
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
