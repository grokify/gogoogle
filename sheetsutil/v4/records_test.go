package sheetsutil

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	driveutil "github.com/grokify/gogoogle/driveutil/v3"
)

type fakeSheets struct {
	mu   sync.Mutex
	reqs []*url.URL
}

func (f *fakeSheets) service(t *testing.T, handler http.HandlerFunc) *sheets.Service {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.reqs = append(f.reqs, r.URL)
		f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	svc, err := sheets.NewService(context.Background(),
		option.WithHTTPClient(srv.Client()), option.WithEndpoint(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func writeJSON(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := w.Write([]byte(body)); err != nil {
		t.Errorf("write response: %v", err)
	}
}

func TestReadRecordsFirstSheet(t *testing.T) {
	f := &fakeSheets{}
	svc := f.service(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v4/spreadsheets/ss1" {
			writeJSON(t, w, `{"sheets":[{"properties":{"title":"Round's 1"}},{"properties":{"title":"Other"}}]}`)
			return
		}
		writeJSON(t, w, `{"values":[
			["id","response","","note"],
			["el-1","agree","x","fine"],
			["el-2","change"],
			[],
			["","","",""],
			["el-3","remove","y",5]
		]}`)
	})
	got, err := ReadRecords(context.Background(), svc, "ss1", "")
	if err != nil {
		t.Fatal(err)
	}
	want := []map[string]string{
		{"id": "el-1", "response": "agree", "note": "fine"},
		{"id": "el-2", "response": "change", "note": ""},
		{"id": "el-3", "response": "remove", "note": "5"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d records, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		for k, v := range want[i] {
			if got[i][k] != v {
				t.Errorf("record %d %q = %q, want %q", i, k, got[i][k], v)
			}
		}
		if len(got[i]) != len(want[i]) {
			t.Errorf("record %d has keys %v (blank-header column should be ignored)", i, got[i])
		}
	}
	if len(f.reqs) != 2 {
		t.Fatalf("want 2 requests, got %d", len(f.reqs))
	}
	if f.reqs[1].Query().Get("valueRenderOption") != "FORMATTED_VALUE" {
		t.Errorf("query = %v", f.reqs[1].RawQuery)
	}
	if p := f.reqs[1].EscapedPath(); p != "/v4/spreadsheets/ss1/values/%27Round%27%27s%201%27" {
		t.Errorf("path = %q: sheet name must be quoted with doubled quotes", p)
	}
}

func TestReadRecordsNamedSheetSkipsMetadataCall(t *testing.T) {
	f := &fakeSheets{}
	svc := f.service(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(t, w, `{"values":[["a"],["1"]]}`) })
	got, err := ReadRecords(context.Background(), svc, "ss1", "Data")
	if err != nil || len(got) != 1 || got[0]["a"] != "1" {
		t.Fatalf("got %v, %v", got, err)
	}
	if len(f.reqs) != 1 {
		t.Errorf("want 1 request, got %d", len(f.reqs))
	}
}

func TestReadRecordsEmptySheet(t *testing.T) {
	svc := (&fakeSheets{}).service(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(t, w, `{}`) })
	got, err := ReadRecords(context.Background(), svc, "ss1", "Data")
	if err != nil || got == nil || len(got) != 0 {
		t.Errorf("want empty non-nil slice, got %v, %v", got, err)
	}
}

func TestReadRecordsErrors(t *testing.T) {
	t.Run("duplicate header", func(t *testing.T) {
		svc := (&fakeSheets{}).service(t, func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, `{"values":[["a","a"],["1","2"]]}`)
		})
		if _, err := ReadRecords(context.Background(), svc, "ss1", "Data"); !errors.Is(err, ErrDuplicateHeader) {
			t.Errorf("err = %v, want ErrDuplicateHeader", err)
		}
	})
	t.Run("no sheets", func(t *testing.T) {
		svc := (&fakeSheets{}).service(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(t, w, `{"sheets":[]}`) })
		if _, err := ReadRecords(context.Background(), svc, "ss1", ""); !errors.Is(err, ErrNoSheets) {
			t.Errorf("err = %v, want ErrNoSheets", err)
		}
	})
	t.Run("api error names the spreadsheet", func(t *testing.T) {
		svc := (&fakeSheets{}).service(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			writeJSON(t, w, `{"error":{"code":404,"message":"gone"}}`)
		})
		_, err := ReadRecords(context.Background(), svc, "ss1", "Data")
		if err == nil || !strings.Contains(err.Error(), "ss1") || !strings.Contains(err.Error(), "Data") {
			t.Errorf("err = %v, want spreadsheet and sheet in message", err)
		}
	})
}

func TestCreateFromCSVConvertsToSheet(t *testing.T) {
	var gotMime string
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, `{"id":"sheet1","name":"data","mimeType":"`+driveutil.MimeGoogleSheet+`","modifiedTime":"2026-01-02T03:04:05Z"}`)
		gotMime = r.Header.Get("Content-Type")
	}))
	t.Cleanup(srv.Close)
	drv, err := driveutil.NewService(context.Background(), srv.Client(), option.WithEndpoint(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "data.csv")
	if err := os.WriteFile(p, []byte("a,b\n1,2\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	f, err := CreateFromCSV(context.Background(), drv, p, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if f.ID != "sheet1" || gotPath != "/upload/drive/v3/files" || gotMime == "" {
		t.Errorf("file %+v, path %q, content-type %q", f, gotPath, gotMime)
	}

	bad := filepath.Join(t.TempDir(), "data.docx")
	if err := os.WriteFile(bad, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateFromCSV(context.Background(), drv, bad, "", ""); !errors.Is(err, driveutil.ErrUnsupportedConversion) {
		t.Errorf("err = %v, want ErrUnsupportedConversion", err)
	}
}
