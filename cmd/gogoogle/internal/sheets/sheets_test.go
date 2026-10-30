package sheets

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/api/option"
	gsheets "google.golang.org/api/sheets/v4"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/cmdutil"
)

func serve(t *testing.T) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"values":[["id","response"],["el-1","agree"],["el-2","change"]]}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	t.Cleanup(srv.Close)
	prev := open
	t.Cleanup(func() { open = prev })
	open = func(ctx context.Context) (*gsheets.Service, error) {
		return gsheets.NewService(ctx, option.WithHTTPClient(srv.Client()), option.WithEndpoint(srv.URL))
	}
}

func TestRecordsPrintsHeaderKeyedJSON(t *testing.T) {
	serve(t)
	cmd := recordsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"https://docs.google.com/spreadsheets/d/ss1/edit", "--sheet", "Round 1"})
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	var got []map[string]string
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out.String())
	}
	if len(got) != 2 || got[0]["id"] != "el-1" || got[1]["response"] != "change" {
		t.Errorf("records = %v", got)
	}
}

func TestRecordsUsageErrors(t *testing.T) {
	serve(t)
	for name, args := range map[string][]string{"no argument": nil, "two arguments": {"a", "b"}, "bad URL": {"https://example.org/x"}} {
		cmd := recordsCmd()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs(args)
		if err := cmd.ExecuteContext(context.Background()); err == nil || !cmdutil.IsUsage(err) {
			t.Errorf("%s: err = %v, want a usage error", name, err)
		}
	}
}
