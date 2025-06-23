package iwark

import (
	"github.com/Iwark/spreadsheet"
	"github.com/grokify/gocharts/v2/data/table"
)

// ParseTableSetFromSpreadsheet is used to parse a TableSet from a sheet. A constraint of this
// function is that `colsRowIndex` and `headerRowCount` are the same value for all sheets.
func ParseTableSetFromSpreadsheet(ss spreadsheet.Spreadsheet, opts *ReadSpreadsheetOpts) (*table.TableSet, error) {
	ts := table.NewTableSet("")
	if opts == nil {
		opts = DefaultReadSheetOpts()
	}
	for _, sheet := range ss.Sheets {
		if !opts.InclHidden && sheet.Properties.Hidden {
			continue
		}
		if t, err := ParseTableFromSheet(&sheet, opts); err != nil {
			return nil, err
		} else if err := ts.Add(t); err != nil {
			return nil, err
		}
	}
	return ts, nil
}
