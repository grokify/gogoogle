// Package sheetsutil provides utilities for working with the Google Sheets API v4.
//
// This package includes:
//
//   - URL parsing for Google Sheets URLs (extracting spreadsheet IDs, sheet GIDs, ranges)
//   - Cell value types for representing typed spreadsheet data
//   - Grid parsing utilities for converting API responses to Go-friendly formats
//
// # URL Parsing
//
// Parse Google Sheets URLs to extract spreadsheet IDs and other metadata:
//
//	id, err := sheetsutil.ParseSpreadsheetURL("https://docs.google.com/spreadsheets/d/abc123/edit#gid=0")
//	// id = "abc123"
//
//	info, err := sheetsutil.ParseSpreadsheetURLFull("https://docs.google.com/spreadsheets/d/abc123/edit#gid=456&range=A1:D10")
//	// info.SpreadsheetID = "abc123"
//	// info.SheetGID = 456
//	// info.Range = "A1:D10"
//
// Build URLs from IDs:
//
//	url := sheetsutil.BuildSpreadsheetURL("abc123")
//	// url = "https://docs.google.com/spreadsheets/d/abc123/edit"
//
//	url := sheetsutil.BuildSpreadsheetURLWithSheet("abc123", 456)
//	// url = "https://docs.google.com/spreadsheets/d/abc123/edit#gid=456"
//
// # Cell Value Types
//
// CellValue provides a JSON-friendly representation of spreadsheet cell data:
//
//	cv := sheetsutil.ParseCellValue(123.45, "$123.45")
//	// cv.Type = CellTypeNumber
//	// cv.NumberValue = 123.45
//	// cv.FormattedValue = "$123.45"
//
// TypedCellValue extends CellValue with Go-native types:
//
//	opts := sheetsutil.ValueParseOptions{PreferInt64: true}
//	tcv := sheetsutil.ParseTypedCellValue(42.0, "42", opts)
//	// tcv.Int64 = 42
//
// # Grid Parsing
//
// Convert Sheets API responses to typed grids:
//
//	vr := &sheets.ValueRange{...}
//	grid := sheetsutil.ParseValueRange(vr)
//	// grid is [][]CellValue
//
// Or extract just the formatted values:
//
//	values := sheetsutil.ExtractFormattedValues(vr)
//	// values is [][]string
package sheetsutil
