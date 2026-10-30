package docs

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/cmdutil"
	docsutil "github.com/grokify/gogoogle/docsutil/v1"
)

const docJSON = `{"documentId":"d1","body":{"content":[{"paragraph":{"elements":[
  {"textRun":{"content":"Keep "}},
  {"textRun":{"content":"old ","suggestedDeletionIds":["s.1"]}},
  {"textRun":{"content":"new ","suggestedInsertionIds":["s.2"]}},
  {"textRun":{"content":"text\n"}}]}}]}}`

func serve(t *testing.T, status int, body string) *string {
	t.Helper()
	var query string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	t.Cleanup(srv.Close)
	prev := open
	t.Cleanup(func() { open = prev })
	open = func(ctx context.Context) (*docsutil.Service, error) {
		svc, err := docs.NewService(ctx, option.WithHTTPClient(srv.Client()), option.WithEndpoint(srv.URL))
		if err != nil {
			return nil, err
		}
		return &docsutil.Service{DocsService: svc, DocumentsService: docs.NewDocumentsService(svc)}, nil
	}
	return &query
}

func TestSuggestionsPrintsClassifiedSegments(t *testing.T) {
	query := serve(t, http.StatusOK, docJSON)
	cmd := suggestionsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"https://docs.google.com/document/d/d1/edit"})
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	var got docsutil.SuggestionContent
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out.String())
	}
	if len(got.Blocks) != 1 || got.Blocks[0].Original != "Keep old text\n" || got.Blocks[0].Accepted != "Keep new text\n" {
		t.Errorf("blocks = %+v", got.Blocks)
	}
	if !strings.Contains(*query, "suggestionsViewMode=SUGGESTIONS_INLINE") {
		t.Errorf("query = %q, want SUGGESTIONS_INLINE", *query)
	}
}

func TestSuggestionsUsageErrors(t *testing.T) {
	serve(t, http.StatusOK, docJSON)
	for name, args := range map[string][]string{"no argument": nil, "two arguments": {"a", "b"}, "bad URL": {"https://example.org/x"}} {
		cmd := suggestionsCmd()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs(args)
		if err := cmd.ExecuteContext(context.Background()); err == nil || !cmdutil.IsUsage(err) {
			t.Errorf("%s: err = %v, want a usage error", name, err)
		}
	}
}

func TestSuggestionsAPIFailure(t *testing.T) {
	serve(t, http.StatusForbidden, `{"error":{"code":403,"message":"no access"}}`)
	cmd := suggestionsCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"d1"})
	err := cmd.ExecuteContext(context.Background())
	if err == nil || cmdutil.IsUsage(err) || !strings.Contains(err.Error(), "d1") {
		t.Errorf("err = %v, want a non-usage error naming the document", err)
	}
}
