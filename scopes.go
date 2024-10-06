package gogoogle

import (
	sheets "google.golang.org/api/sheets/v4"
)

const (
	ScopeDrive          = sheets.DriveScope
	ScopeDriveFile      = sheets.DriveFileScope
	ScopeDriveReadonly  = sheets.DriveReadonlyScope
	ScopeSlides         = sheets.SpreadsheetsScope
	ScopeSlidesReadonly = sheets.SpreadsheetsReadonlyScope
)
