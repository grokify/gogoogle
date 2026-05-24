package docsutil

import (
	"testing"
)

func TestParseDocumentURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "plain ID",
			input: "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			want:  "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
		},
		{
			name:  "basic URL",
			input: "https://docs.google.com/document/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			want:  "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
		},
		{
			name:  "URL with /edit",
			input: "https://docs.google.com/document/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit",
			want:  "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
		},
		{
			name:  "URL with query string",
			input: "https://docs.google.com/document/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit?tab=t.0",
			want:  "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
		},
		{
			name:  "URL with fragment",
			input: "https://docs.google.com/document/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit#heading=h.nyn35zfb1x8z",
			want:  "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
		},
		{
			name:  "URL with query string and fragment",
			input: "https://docs.google.com/document/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit?tab=t.0#heading=h.nyn35zfb1x8z",
			want:  "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
		},
		{
			name:  "HTTP URL",
			input: "http://docs.google.com/document/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit",
			want:  "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
		},
		{
			name:  "ID with underscores and dashes",
			input: "abc_DEF-123_xyz",
			want:  "abc_DEF-123_xyz",
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "invalid URL - wrong domain",
			input:   "https://google.com/document/d/123",
			wantErr: true,
		},
		{
			name:    "invalid URL - spreadsheet",
			input:   "https://docs.google.com/spreadsheets/d/123/edit",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDocumentURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDocumentURL(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseDocumentURL(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDocumentURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildDocumentURL(t *testing.T) {
	id := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
	want := "https://docs.google.com/document/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit"

	got := BuildDocumentURL(id)
	if got != want {
		t.Errorf("BuildDocumentURL(%q) = %q, want %q", id, got, want)
	}
}

func TestIsDocumentURL(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"https://docs.google.com/document/d/abc123/edit", true},
		{"https://docs.google.com/document/d/abc123/edit?tab=t.0#heading=h.xyz", true},
		{"https://docs.google.com/spreadsheets/d/abc123/edit", false},
		{"abc123", false},
		{"https://google.com/document/d/abc123", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := IsDocumentURL(tt.input)
			if got != tt.want {
				t.Errorf("IsDocumentURL(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
