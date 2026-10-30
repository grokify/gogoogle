package sheetsutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/api/option"
)

func TestServiceGetSpreadsheet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/v4/spreadsheets/ss1" {
			w.WriteHeader(http.StatusNotFound)
			writeJSON(t, w, `{"error":{"code":404,"message":"gone"}}`)
			return
		}
		writeJSON(t, w, `{"spreadsheetId":"ss1","properties":{"title":"Round 1"}}`)
	}))
	t.Cleanup(srv.Close)
	svc, err := NewService(context.Background(), srv.Client(), option.WithEndpoint(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	ss, err := svc.GetSpreadsheet(context.Background(), "ss1")
	if err != nil || ss.Properties.Title != "Round 1" {
		t.Fatalf("got %+v, %v", ss, err)
	}
	if _, err := svc.GetSpreadsheet(context.Background(), "missing"); err == nil || !strings.Contains(err.Error(), "missing") {
		t.Errorf("err = %v, want one naming the spreadsheet", err)
	}
}
