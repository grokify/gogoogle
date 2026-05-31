package sheetsutil

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrInvalidURL is returned when a URL cannot be parsed.
	ErrInvalidURL = errors.New("invalid Google Sheets URL")

	// ErrEmptyInput is returned when input is empty.
	ErrEmptyInput = errors.New("input cannot be empty")

	// rxSheetsURL matches Google Sheets spreadsheet URLs.
	// Supports formats:
	//   - https://docs.google.com/spreadsheets/d/{id}
	//   - https://docs.google.com/spreadsheets/d/{id}/edit
	//   - https://docs.google.com/spreadsheets/d/{id}/edit#gid={sheet_gid}
	//   - https://docs.google.com/spreadsheets/d/{id}/edit?gid={sheet_gid}
	//   - https://docs.google.com/spreadsheets/d/{id}/edit#gid={sheet_gid}&range=A1:D10
	// The spreadsheet ID is captured in group 1.
	rxSheetsURL = regexp.MustCompile(`(?i)^https?://docs\.google\.com/spreadsheets/d/([a-zA-Z0-9_-]+)(?:/[^?#]*)?(?:\?[^#]*)?(?:#.*)?$`)
)

// SpreadsheetURLInfo contains parsed information from a Google Sheets URL.
type SpreadsheetURLInfo struct {
	SpreadsheetID string `json:"spreadsheet_id"`
	SheetGID      *int64 `json:"sheet_gid,omitempty"`
	Range         string `json:"range,omitempty"`
}

// ParseSpreadsheetURL extracts the spreadsheet ID from a Google Sheets URL.
// If the input is already a plain ID (no slashes or protocol), it returns it as-is.
// Returns the spreadsheet ID or an error if the URL is invalid.
func ParseSpreadsheetURL(urlOrID string) (string, error) {
	urlOrID = strings.TrimSpace(urlOrID)
	if urlOrID == "" {
		return "", ErrEmptyInput
	}

	// Check if it's already just an ID (no slashes or protocol)
	if !strings.Contains(urlOrID, "/") && !strings.Contains(urlOrID, ":") {
		return urlOrID, nil
	}

	// Try to parse as URL
	matches := rxSheetsURL.FindStringSubmatch(urlOrID)
	if len(matches) < 2 {
		return "", ErrInvalidURL
	}

	return matches[1], nil
}

// ParseSpreadsheetURLFull extracts all available info from a Google Sheets URL.
// This includes the spreadsheet ID, optional sheet GID, and optional range.
func ParseSpreadsheetURLFull(urlOrID string) (SpreadsheetURLInfo, error) {
	info := SpreadsheetURLInfo{}

	urlOrID = strings.TrimSpace(urlOrID)
	if urlOrID == "" {
		return info, ErrEmptyInput
	}

	// Check if it's already just an ID (no slashes or protocol)
	if !strings.Contains(urlOrID, "/") && !strings.Contains(urlOrID, ":") {
		info.SpreadsheetID = urlOrID
		return info, nil
	}

	// Try to parse as URL
	matches := rxSheetsURL.FindStringSubmatch(urlOrID)
	if len(matches) < 2 {
		return info, ErrInvalidURL
	}
	info.SpreadsheetID = matches[1]

	// Parse the URL for query params and fragment
	parsedURL, err := url.Parse(urlOrID)
	if err != nil {
		// We already got the ID, so return what we have
		return info, nil
	}

	// Try to get gid from query params (newer format: ?gid=123)
	if gidStr := parsedURL.Query().Get("gid"); gidStr != "" {
		if gid, err := strconv.ParseInt(gidStr, 10, 64); err == nil {
			info.SheetGID = &gid
		}
	}

	// Parse fragment for gid and range (older format: #gid=123&range=A1:D10)
	if parsedURL.Fragment != "" {
		fragmentParams := parseFragment(parsedURL.Fragment)

		// Get gid from fragment if not already set from query
		if info.SheetGID == nil {
			if gidStr, ok := fragmentParams["gid"]; ok {
				if gid, err := strconv.ParseInt(gidStr, 10, 64); err == nil {
					info.SheetGID = &gid
				}
			}
		}

		// Get range from fragment
		if rangeStr, ok := fragmentParams["range"]; ok {
			info.Range = rangeStr
		}
	}

	return info, nil
}

// parseFragment parses URL fragment parameters (e.g., "gid=123&range=A1:D10")
func parseFragment(fragment string) map[string]string {
	params := make(map[string]string)
	pairs := strings.Split(fragment, "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}
	return params
}

// BuildSpreadsheetURL constructs a Google Sheets URL from a spreadsheet ID.
func BuildSpreadsheetURL(spreadsheetID string) string {
	return "https://docs.google.com/spreadsheets/d/" + spreadsheetID + "/edit"
}

// BuildSpreadsheetURLWithSheet constructs a Google Sheets URL with a sheet GID.
func BuildSpreadsheetURLWithSheet(spreadsheetID string, sheetGID int64) string {
	return "https://docs.google.com/spreadsheets/d/" + spreadsheetID + "/edit#gid=" + strconv.FormatInt(sheetGID, 10)
}

// IsSpreadsheetURL checks if the given string is a Google Sheets spreadsheet URL.
func IsSpreadsheetURL(s string) bool {
	return rxSheetsURL.MatchString(s)
}

// NormalizeSpreadsheetInput accepts either a spreadsheet ID or URL and returns the spreadsheet ID.
// This is useful for APIs that accept either format.
func NormalizeSpreadsheetInput(urlOrID string) (string, error) {
	return ParseSpreadsheetURL(urlOrID)
}
