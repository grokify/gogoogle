package sheetsutil

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/sheets/v4"

	driveutil "github.com/grokify/gogoogle/driveutil/v3"
)

// ErrNoSheets is returned when a spreadsheet has no sheets to read.
var ErrNoSheets = errors.New("spreadsheet has no sheets")

// ErrDuplicateHeader is returned when two columns share a header.
var ErrDuplicateHeader = errors.New("duplicate header")

// ReadRecords reads one sheet as records keyed by the header row (row 1).
// sheetName selects the sheet; empty means the first sheet.
//
// Values are the formatted strings a person sees. Columns with a blank header
// are ignored. A short row is padded with empty strings, and a row with no
// values is skipped. A repeated header is an error, since one value would
// silently shadow the other. The returned slice is in sheet order.
//
// Reading needs the spreadsheets.readonly scope or a broader one.
func ReadRecords(ctx context.Context, svc *sheets.Service, spreadsheetID, sheetName string) ([]map[string]string, error) {
	if sheetName == "" {
		ss, err := svc.Spreadsheets.Get(spreadsheetID).Fields("sheets.properties.title").Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("sheetsutil: get spreadsheet %s: %w", spreadsheetID, err)
		}
		if len(ss.Sheets) == 0 || ss.Sheets[0].Properties == nil {
			return nil, fmt.Errorf("sheetsutil: read records %s: %w", spreadsheetID, ErrNoSheets)
		}
		sheetName = ss.Sheets[0].Properties.Title
	}
	vr, err := svc.Spreadsheets.Values.Get(spreadsheetID, quoteSheetName(sheetName)).
		ValueRenderOption("FORMATTED_VALUE").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("sheetsutil: read values %s %q: %w", spreadsheetID, sheetName, err)
	}
	records, err := recordsFromRows(vr.Values)
	if err != nil {
		return nil, fmt.Errorf("sheetsutil: read records %s %q: %w", spreadsheetID, sheetName, err)
	}
	return records, nil
}

// recordsFromRows turns a header row and data rows into records.
func recordsFromRows(rows [][]any) ([]map[string]string, error) {
	records := []map[string]string{}
	if len(rows) == 0 {
		return records, nil
	}
	headers := make([]string, len(rows[0]))
	seen := map[string]bool{}
	for i, h := range rows[0] {
		name := strings.TrimSpace(cellString(h))
		if name == "" {
			continue
		}
		if seen[name] {
			return nil, fmt.Errorf("%w %q", ErrDuplicateHeader, name)
		}
		seen[name] = true
		headers[i] = name
	}
	for _, row := range rows[1:] {
		rec := make(map[string]string, len(seen))
		empty := true
		for i, name := range headers {
			if name == "" {
				continue
			}
			v := ""
			if i < len(row) {
				v = cellString(row[i])
			}
			if v != "" {
				empty = false
			}
			rec[name] = v
		}
		if !empty {
			records = append(records, rec)
		}
	}
	return records, nil
}

func cellString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

// quoteSheetName wraps a sheet name for A1 notation, doubling single quotes.
func quoteSheetName(name string) string {
	return "'" + strings.ReplaceAll(name, "'", "''") + "'"
}

// CreateFromCSV uploads a CSV file as a new Google Sheet and returns the
// created file. name defaults to the file's base name; folderID is optional.
// It delegates to the Drive upload with conversion, so it needs the
// drive.file scope.
func CreateFromCSV(ctx context.Context, drv *driveutil.Service, path, name, folderID string) (*driveutil.File, error) {
	f, err := drv.Upload(ctx, driveutil.UploadRequest{
		Path:      path,
		Name:      name,
		FolderID:  folderID,
		ConvertTo: driveutil.ConvertSheet,
	})
	if err != nil {
		return nil, fmt.Errorf("sheetsutil: create from csv %s: %w", path, err)
	}
	return f, nil
}
