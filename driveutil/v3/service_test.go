package driveutil

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type recorded struct {
	Method string
	Path   string
	Query  map[string][]string
	CType  string
	Body   string
}

type fakeServer struct {
	*httptest.Server
	mu   sync.Mutex
	reqs []recorded
}

func newFake(t *testing.T, h func(w http.ResponseWriter, r *http.Request, n int)) *fakeServer {
	t.Helper()
	f := &fakeServer{}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		f.mu.Lock()
		f.reqs = append(f.reqs, recorded{r.Method, r.URL.Path, r.URL.Query(), r.Header.Get("Content-Type"), string(buf)})
		n := len(f.reqs)
		f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		h(w, r, n)
	}))
	t.Cleanup(f.Close)
	return f
}

func (f *fakeServer) requests() []recorded {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]recorded(nil), f.reqs...)
}

func (f *fakeServer) service(t *testing.T) *Service {
	t.Helper()
	svc, err := NewService(context.Background(), f.Client(), option.WithEndpoint(f.URL+"/"))
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestGetFile(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeBody(t, w, `{"id":"f1","name":"Doc","mimeType":"application/vnd.google-apps.document","webViewLink":"https://x/f1","modifiedTime":"2026-01-02T03:04:05.000Z","headRevisionId":"rev9"}`)
	})
	got, err := f.service(t).GetFile(context.Background(), "f1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "f1" || got.Name != "Doc" || got.WebViewLink != "https://x/f1" || got.HeadRevisionID != "rev9" {
		t.Fatalf("unexpected file: %+v", got)
	}
	if got.ModifiedTime.Year() != 2026 || got.ModifiedTime.Minute() != 4 {
		t.Fatalf("bad modified time: %v", got.ModifiedTime)
	}
	reqs := f.requests()
	if len(reqs) != 1 || reqs[0].Method != "GET" || reqs[0].Path != "/files/f1" {
		t.Fatalf("unexpected requests: %+v", reqs)
	}
	q := reqs[0].Query
	if q["fields"][0] != fileFields || q["supportsAllDrives"][0] != "true" {
		t.Fatalf("unexpected query: %v", q)
	}
}

func TestGetFileBadTime(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeBody(t, w, `{"id":"f1","modifiedTime":"yesterday"}`)
	})
	if _, err := f.service(t).GetFile(context.Background(), "f1"); err == nil {
		t.Fatal("expected error for bad modifiedTime")
	}
}

func TestGetFileErrors(t *testing.T) {
	tests := []struct {
		code int
		want error
	}{
		{http.StatusNotFound, ErrNotFound},
		{http.StatusForbidden, ErrPermissionDenied},
	}
	for _, tc := range tests {
		f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
			w.WriteHeader(tc.code)
			writeBody(t, w, `{"error":{"code":`+strconv.Itoa(tc.code)+`,"message":"nope"}}`)
		})
		_, err := f.service(t).GetFile(context.Background(), "zz")
		if !errors.Is(err, tc.want) {
			t.Fatalf("code %d: want %v, got %v", tc.code, tc.want, err)
		}
		var gerr *googleapi.Error
		if !errors.As(err, &gerr) || gerr.Code != tc.code {
			t.Fatalf("code %d: googleapi.Error not wrapped: %v", tc.code, err)
		}
	}
}

func TestGetFileOtherError(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(http.StatusInternalServerError)
		writeBody(t, w, `{"error":{"code":500,"message":"boom"}}`)
	})
	_, err := f.service(t).GetFile(context.Background(), "zz")
	var gerr *googleapi.Error
	if !errors.As(err, &gerr) || gerr.Code != 500 {
		t.Fatalf("expected wrapped 500, got %v", err)
	}
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("500 must not match sentinels: %v", err)
	}
}

func TestWrapErrNil(t *testing.T) {
	if wrapErr("x", "y", nil) != nil {
		t.Fatal("nil must stay nil")
	}
}

func writeBody(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := w.Write([]byte(body)); err != nil {
		t.Errorf("write response: %v", err)
	}
}
