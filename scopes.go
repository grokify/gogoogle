package gogoogle

import (
	sheets "google.golang.org/api/sheets/v4"
)

const (
	ScopeDrive                = sheets.DriveScope
	ScopeDriveFile            = sheets.DriveFileScope
	ScopeDriveReadonly        = sheets.DriveReadonlyScope
	ScopeSpreadsheets         = sheets.SpreadsheetsScope
	ScopeSpreadsheetsReadonly = sheets.SpreadsheetsReadonlyScope // "https://www.googleapis.com/auth/spreadsheets.readonly"
)
