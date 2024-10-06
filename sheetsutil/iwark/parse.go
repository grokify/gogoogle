package iwark

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Iwark/spreadsheet"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gogoogle/docsutil"
	"github.com/grokify/mogo/type/stringsutil"
)

var ErrSheetIDRequired = errors.New("sheet id is required")

func ParseTableFromSheetID(client *http.Client, sheetID string, sheetIdx, headerRows uint) (*table.Table, error) {
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

func ParseTableFromSpreadsheet(ss spreadsheet.Spreadsheet, sheetIdx, headerRows uint) (*table.Table, error) {
	if s, err := ss.SheetByIndex(sheetIdx); err != nil {
		return nil, err
	} else {
		return ParseTableFromSheet(s, headerRows), nil
	}
}

func ParseTableFromSheet(s *spreadsheet.Sheet, headerRows uint) *table.Table {
	cols, rows := ParseDataFromSheet(s, headerRows)
	tbl := table.NewTable("")
	tbl.Columns = cols
	tbl.Rows = rows
	return &tbl
}

func ParseDataFromSheet(s *spreadsheet.Sheet, headerRows uint) ([]string, [][]string) {
	var cols []string
	var rows [][]string
	for i, srow := range s.Rows {
		var row []string
		for _, scell := range srow {
			row = append(row, scell.Value)
		}
		if headerRows > 0 && uint(i) < headerRows {
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
