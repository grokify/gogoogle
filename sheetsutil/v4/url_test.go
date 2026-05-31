package sheetsutil

import (
	"testing"
)

// Test cases use Google's official public sample spreadsheet ID:
// "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
// This spreadsheet contains sample class roster data and is used in Google's
// Sheets API documentation and tutorials. It is publicly accessible at:
// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms

func TestParseSpreadsheetURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "plain ID",
			input:   "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
		{
			name:    "basic URL",
			input:   "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
		{
			name:    "URL with /edit",
			input:   "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
		{
			name:    "URL with gid fragment",
			input:   "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit#gid=123",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
		{
			name:    "URL with gid query param",
			input:   "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit?gid=456",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
		{
			name:    "URL with gid and range",
			input:   "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit#gid=789&range=A1:D10",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
		{
			name:    "URL with http",
			input:   "http://docs.google.com/spreadsheets/d/abc123_-XYZ/edit",
			want:    "abc123_-XYZ",
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid URL - wrong domain",
			input:   "https://example.com/spreadsheets/d/abc123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid URL - docs URL but wrong type",
			input:   "https://docs.google.com/document/d/abc123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "ID with leading/trailing whitespace",
			input:   "  1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms  ",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSpreadsheetURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSpreadsheetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseSpreadsheetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSpreadsheetURLFull(t *testing.T) {
	gid123 := int64(123)
	gid456 := int64(456)
	gid789 := int64(789)
	gid0 := int64(0)

	tests := []struct {
		name    string
		input   string
		wantID  string
		wantGID *int64
		wantRng string
		wantErr bool
	}{
		{
			name:    "plain ID",
			input:   "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantID:  "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantGID: nil,
			wantRng: "",
			wantErr: false,
		},
		{
			name:    "URL with gid fragment",
			input:   "https://docs.google.com/spreadsheets/d/abc123/edit#gid=123",
			wantID:  "abc123",
			wantGID: &gid123,
			wantRng: "",
			wantErr: false,
		},
		{
			name:    "URL with gid query param",
			input:   "https://docs.google.com/spreadsheets/d/abc123/edit?gid=456",
			wantID:  "abc123",
			wantGID: &gid456,
			wantRng: "",
			wantErr: false,
		},
		{
			name:    "URL with gid=0 (first sheet)",
			input:   "https://docs.google.com/spreadsheets/d/abc123/edit#gid=0",
			wantID:  "abc123",
			wantGID: &gid0,
			wantRng: "",
			wantErr: false,
		},
		{
			name:    "URL with gid and range in fragment",
			input:   "https://docs.google.com/spreadsheets/d/abc123/edit#gid=789&range=A1:D10",
			wantID:  "abc123",
			wantGID: &gid789,
			wantRng: "A1:D10",
			wantErr: false,
		},
		{
			name:    "URL with range only in fragment",
			input:   "https://docs.google.com/spreadsheets/d/abc123/edit#range=B2:E5",
			wantID:  "abc123",
			wantGID: nil,
			wantRng: "B2:E5",
			wantErr: false,
		},
		{
			name:    "URL without any params",
			input:   "https://docs.google.com/spreadsheets/d/abc123/edit",
			wantID:  "abc123",
			wantGID: nil,
			wantRng: "",
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			wantID:  "",
			wantGID: nil,
			wantRng: "",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			input:   "https://example.com/sheets/abc",
			wantID:  "",
			wantGID: nil,
			wantRng: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSpreadsheetURLFull(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSpreadsheetURLFull() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.SpreadsheetID != tt.wantID {
				t.Errorf("ParseSpreadsheetURLFull().SpreadsheetID = %v, want %v", got.SpreadsheetID, tt.wantID)
			}
			if (got.SheetGID == nil) != (tt.wantGID == nil) {
				t.Errorf("ParseSpreadsheetURLFull().SheetGID = %v, want %v", got.SheetGID, tt.wantGID)
			} else if got.SheetGID != nil && *got.SheetGID != *tt.wantGID {
				t.Errorf("ParseSpreadsheetURLFull().SheetGID = %v, want %v", *got.SheetGID, *tt.wantGID)
			}
			if got.Range != tt.wantRng {
				t.Errorf("ParseSpreadsheetURLFull().Range = %v, want %v", got.Range, tt.wantRng)
			}
		})
	}
}

func TestBuildSpreadsheetURL(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "basic ID",
			id:   "abc123",
			want: "https://docs.google.com/spreadsheets/d/abc123/edit",
		},
		{
			name: "long ID",
			id:   "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			want: "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSpreadsheetURL(tt.id)
			if got != tt.want {
				t.Errorf("BuildSpreadsheetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildSpreadsheetURLWithSheet(t *testing.T) {
	tests := []struct {
		name string
		id   string
		gid  int64
		want string
	}{
		{
			name: "basic with gid=0",
			id:   "abc123",
			gid:  0,
			want: "https://docs.google.com/spreadsheets/d/abc123/edit#gid=0",
		},
		{
			name: "with gid=123456",
			id:   "abc123",
			gid:  123456,
			want: "https://docs.google.com/spreadsheets/d/abc123/edit#gid=123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSpreadsheetURLWithSheet(tt.id, tt.gid)
			if got != tt.want {
				t.Errorf("BuildSpreadsheetURLWithSheet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSpreadsheetURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid sheets URL",
			input: "https://docs.google.com/spreadsheets/d/abc123/edit",
			want:  true,
		},
		{
			name:  "valid sheets URL with gid",
			input: "https://docs.google.com/spreadsheets/d/abc123/edit#gid=0",
			want:  true,
		},
		{
			name:  "plain ID",
			input: "abc123",
			want:  false,
		},
		{
			name:  "docs URL",
			input: "https://docs.google.com/document/d/abc123/edit",
			want:  false,
		},
		{
			name:  "slides URL",
			input: "https://docs.google.com/presentation/d/abc123/edit",
			want:  false,
		},
		{
			name:  "random URL",
			input: "https://example.com/abc123",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSpreadsheetURL(tt.input)
			if got != tt.want {
				t.Errorf("IsSpreadsheetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeSpreadsheetInput(t *testing.T) {
	// NormalizeSpreadsheetInput is an alias for ParseSpreadsheetURL, so just test basic functionality
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "ID passthrough",
			input:   "abc123",
			want:    "abc123",
			wantErr: false,
		},
		{
			name:    "URL extraction",
			input:   "https://docs.google.com/spreadsheets/d/xyz789/edit",
			want:    "xyz789",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeSpreadsheetInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeSpreadsheetInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeSpreadsheetInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
