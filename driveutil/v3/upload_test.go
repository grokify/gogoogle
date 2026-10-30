package driveutil

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

const createdJSON = `{"id":"new1","name":"report","mimeType":"x","modifiedTime":"2026-01-02T03:04:05Z"}`

func TestUploadConversions(t *testing.T) {
	tests := []struct {
		file     string
		to       ConvertTo
		wantMeta string
		wantName string
		wantType string
	}{
		{"report.docx", ConvertDoc, MimeGoogleDoc, "report", "wordprocessingml.document"},
		{"notes.txt", ConvertDoc, MimeGoogleDoc, "notes", "text/plain"},
		{"page.HTML", ConvertDoc, MimeGoogleDoc, "page", "text/html"},
		{"a.odt", ConvertDoc, MimeGoogleDoc, "a", "opendocument.text"},
		{"data.csv", ConvertSheet, MimeGoogleSheet, "data", "text/csv"},
		{"data.xlsx", ConvertSheet, MimeGoogleSheet, "data", "spreadsheetml.sheet"},
		{"data.ods", ConvertSheet, MimeGoogleSheet, "data", "opendocument.spreadsheet"},
		{"deck.pptx", ConvertSlides, MimeGoogleSlides, "deck", "presentationml.presentation"},
		{"deck.odp", ConvertSlides, MimeGoogleSlides, "deck", "opendocument.presentation"},
	}
	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
			p := writeTemp(t, tc.file, "PAYLOAD")
			got, err := f.service(t).Upload(context.Background(), UploadRequest{Path: p, ConvertTo: tc.to, FolderID: "folder9"})
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != "new1" {
				t.Fatalf("id = %q", got.ID)
			}
			reqs := f.requests()
			if len(reqs) != 1 {
				t.Fatalf("want 1 request, got %d", len(reqs))
			}
			r := reqs[0]
			if r.Method != "POST" || r.Path != "/upload/drive/v3/files" {
				t.Fatalf("unexpected request %s %s", r.Method, r.Path)
			}
			if r.Query["uploadType"][0] != "multipart" || r.Query["supportsAllDrives"][0] != "true" || r.Query["fields"][0] != fileFields {
				t.Fatalf("unexpected query %v", r.Query)
			}
			for _, want := range []string{`"mimeType":"` + tc.wantMeta + `"`, `"name":"` + tc.wantName + `"`, `"parents":["folder9"]`, tc.wantType, "PAYLOAD"} {
				if !strings.Contains(r.Body, want) {
					t.Errorf("body missing %q:\n%s", want, r.Body)
				}
			}
		})
	}
}

func TestUploadNoConversion(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
	p := writeTemp(t, "blob.unknownext", "X")
	if _, err := f.service(t).Upload(context.Background(), UploadRequest{Path: p}); err != nil {
		t.Fatal(err)
	}
	body := f.requests()[0].Body
	if !strings.Contains(body, `"name":"blob.unknownext"`) || strings.Contains(body, `"mimeType"`) || strings.Contains(body, "parents") {
		t.Fatalf("unexpected body:\n%s", body)
	}
	if !strings.Contains(body, "application/octet-stream") {
		t.Fatalf("expected octet-stream media:\n%s", body)
	}
}

func TestUploadExplicitName(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
	p := writeTemp(t, "a.csv", "x,y")
	if _, err := f.service(t).Upload(context.Background(), UploadRequest{Path: p, Name: "Custom", ConvertTo: ConvertSheet}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.requests()[0].Body, `"name":"Custom"`) {
		t.Fatalf("name not used:\n%s", f.requests()[0].Body)
	}
}

func TestUploadMismatchNoRequest(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
	svc := f.service(t)
	p := writeTemp(t, "data.csv", "x")
	_, err := svc.Upload(context.Background(), UploadRequest{Path: p, ConvertTo: ConvertDoc})
	if !errors.Is(err, ErrUnsupportedConversion) {
		t.Fatalf("got %v", err)
	}
	_, err = svc.Upload(context.Background(), UploadRequest{Path: p, ConvertTo: "bogus"})
	if !errors.Is(err, ErrUnsupportedConversion) {
		t.Fatalf("unknown target: got %v", err)
	}
	_, err = svc.Update(context.Background(), "id1", p, ConvertSlides)
	if !errors.Is(err, ErrUnsupportedConversion) {
		t.Fatalf("update: got %v", err)
	}
	if n := len(f.requests()); n != 0 {
		t.Fatalf("expected no requests, got %d", n)
	}
}

func TestUploadMissingFile(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
	svc := f.service(t)
	missing := filepath.Join(t.TempDir(), "nope.csv")
	if _, err := svc.Upload(context.Background(), UploadRequest{Path: missing}); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("upload: got %v", err)
	}
	if _, err := svc.Update(context.Background(), "id1", missing, ""); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("update: got %v", err)
	}
}

func TestUploadAPIError(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(http.StatusForbidden)
		writeBody(t, w, `{"error":{"code":403,"message":"no"}}`)
	})
	p := writeTemp(t, "a.csv", "x")
	_, err := f.service(t).Upload(context.Background(), UploadRequest{Path: p})
	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("got %v", err)
	}
}

func TestUpdateKeepsID(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeBody(t, w, `{"id":"existing7","name":"n","modifiedTime":"2026-01-02T03:04:05Z"}`)
	})
	p := writeTemp(t, "deck.pptx", "SLIDES")
	got, err := f.service(t).Update(context.Background(), "existing7", p, ConvertSlides)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "existing7" {
		t.Fatalf("id = %q", got.ID)
	}
	r := f.requests()[0]
	if r.Method != "PATCH" || r.Path != "/upload/drive/v3/files/existing7" {
		t.Fatalf("unexpected request %s %s", r.Method, r.Path)
	}
	if r.Query["uploadType"][0] != "multipart" || r.Query["supportsAllDrives"][0] != "true" {
		t.Fatalf("unexpected query %v", r.Query)
	}
	if !strings.Contains(r.Body, MimeGoogleSlides) || !strings.Contains(r.Body, "SLIDES") {
		t.Fatalf("unexpected body:\n%s", r.Body)
	}
}

func TestUpdateNoConversionAndError(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(http.StatusNotFound)
		writeBody(t, w, `{"error":{"code":404,"message":"gone"}}`)
	})
	p := writeTemp(t, "a.txt", "x")
	_, err := f.service(t).Update(context.Background(), "gone1", p, ConvertNone)
	if !errors.Is(err, ErrNotFound) || !strings.Contains(err.Error(), "gone1") {
		t.Fatalf("got %v", err)
	}
	if strings.Contains(f.requests()[0].Body, `"mimeType"`) {
		t.Fatalf("metadata should not set mimeType:\n%s", f.requests()[0].Body)
	}
}
