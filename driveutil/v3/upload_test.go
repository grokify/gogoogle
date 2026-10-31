package driveutil

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"
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

// Content is sent in place of the file at Path, which is only a name: it need
// not exist, and its extension still picks the media type and the conversion.
func TestUploadFromContent(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
	got, err := f.service(t).Upload(context.Background(), UploadRequest{
		Path:      filepath.Join(t.TempDir(), "gone", "report.docx"), // never opened
		Content:   strings.NewReader("FROM READER"),
		ConvertTo: ConvertDoc,
		FolderID:  "folder9",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "new1" {
		t.Fatalf("id = %q", got.ID)
	}
	r := f.requests()[0]
	for _, want := range []string{`"mimeType":"` + MimeGoogleDoc + `"`, `"name":"report"`, `"parents":["folder9"]`, "wordprocessingml.document", "FROM READER"} {
		if !strings.Contains(r.Body, want) {
			t.Errorf("body missing %q:\n%s", want, r.Body)
		}
	}
}

func TestUploadContentWinsOverTheFileAtPath(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
	p := writeTemp(t, "a.txt", "ON DISK")
	if _, err := f.service(t).Upload(context.Background(), UploadRequest{Path: p, Content: strings.NewReader("FROM READER")}); err != nil {
		t.Fatal(err)
	}
	body := f.requests()[0].Body
	if !strings.Contains(body, "FROM READER") || strings.Contains(body, "ON DISK") {
		t.Fatalf("the file at Path was read instead of Content:\n%s", body)
	}
}

// An error from the reader comes back wrapped, so a caller can tell why the
// upload stopped, for instance that its size cap was exceeded.
func TestUploadContentReaderError(t *testing.T) {
	boom := errors.New("reader failed")
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
	content := io.MultiReader(strings.NewReader("abc"), iotest.ErrReader(boom))
	got, err := f.service(t).Upload(context.Background(), UploadRequest{Path: "a.txt", Content: content})
	if !errors.Is(err, boom) || got != nil {
		t.Fatalf("got %v, %v; want the reader's error and no file", got, err)
	}
}

func TestUpdateContent(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeBody(t, w, `{"id":"existing7","name":"n","modifiedTime":"2026-01-02T03:04:05Z"}`)
	})
	got, err := f.service(t).UpdateContent(context.Background(), "existing7", "deck.pptx", strings.NewReader("SLIDES"), ConvertSlides)
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
	if !strings.Contains(r.Body, MimeGoogleSlides) || !strings.Contains(r.Body, "SLIDES") {
		t.Fatalf("unexpected body:\n%s", r.Body)
	}
}

func TestUpdateContentRefusals(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, createdJSON) })
	svc := f.service(t)
	if _, err := svc.UpdateContent(context.Background(), "id1", "a.txt", nil, ConvertNone); err == nil {
		t.Error("nil content was accepted")
	}
	// The name's extension must match the conversion, as it must for a path.
	if _, err := svc.UpdateContent(context.Background(), "id1", "a.txt", strings.NewReader("x"), ConvertSlides); err == nil {
		t.Error("a mismatched name was accepted")
	}
	boom := errors.New("reader failed")
	if _, err := svc.UpdateContent(context.Background(), "id1", "a.txt", io.MultiReader(strings.NewReader("a"), iotest.ErrReader(boom)), ConvertNone); !errors.Is(err, boom) {
		t.Errorf("reader error: got %v", err)
	}
}

// Update with a path is UpdateContent over the opened file.
func TestUpdateStillReadsThePath(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeBody(t, w, `{"id":"e1","name":"n","modifiedTime":"2026-01-02T03:04:05Z"}`)
	})
	if _, err := f.service(t).Update(context.Background(), "e1", writeTemp(t, "a.txt", "ON DISK"), ConvertNone); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.requests()[0].Body, "ON DISK") {
		t.Fatalf("the file was not sent:\n%s", f.requests()[0].Body)
	}
}
