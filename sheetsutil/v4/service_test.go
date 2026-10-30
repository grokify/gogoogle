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

func TestServiceValues(t *testing.T) {
	var queries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.Path+"?"+r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "batchGet") {
			writeJSON(t, w, `{"valueRanges":[{"range":"A1:B1","values":[["a","b"]]},{"range":"C1","values":[["c"]]}]}`)
			return
		}
		writeJSON(t, w, `{"range":"Sheet1!A1:B1","values":[["a","b"]]}`)
	}))
	t.Cleanup(srv.Close)
	svc, err := NewService(context.Background(), srv.Client(), option.WithEndpoint(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	vr, err := svc.GetValues(context.Background(), "ss1", "Sheet1!A1:B1", "UNFORMATTED_VALUE")
	if err != nil || vr.Range != "Sheet1!A1:B1" {
		t.Fatalf("GetValues = %+v, %v", vr, err)
	}
	resp, err := svc.BatchGetValues(context.Background(), "ss1", []string{"A1:B1", "C1"}, "")
	if err != nil || len(resp.ValueRanges) != 2 {
		t.Fatalf("BatchGetValues = %+v, %v", resp, err)
	}
	if !strings.Contains(queries[0], "valueRenderOption=UNFORMATTED_VALUE") {
		t.Errorf("query %q lacks the render option", queries[0])
	}
	if strings.Contains(queries[1], "valueRenderOption") || !strings.Contains(queries[1], "ranges=C1") {
		t.Errorf("query %q: want ranges and no render option", queries[1])
	}
}
