package sheetsutil

import (
	"fmt"
)

const WebURLPattern = `https://docs.google.com/spreadsheets/d/%s/edit#gid=0`

func SheetToWebURL(spreadsheetId string) string {
	return fmt.Sprintf(WebURLPattern, spreadsheetId)
}
