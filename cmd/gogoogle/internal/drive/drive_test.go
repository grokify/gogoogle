package drive

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"google.golang.org/api/option"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/cmdutil"
	driveutil "github.com/grokify/gogoogle/driveutil/v3"
)

// serve points the commands at a fake Drive API and returns the paths and
// scopes it saw.
func serve(t *testing.T, handler http.HandlerFunc) (paths, scopes *[]string) {
	t.Helper()
	paths, scopes = new([]string), new([]string)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*paths = append(*paths, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	prev := open
	t.Cleanup(func() { open = prev })
	open = func(ctx context.Context, s ...string) (*driveutil.Service, error) {
		*scopes = append(*scopes, s...)
		return driveutil.NewService(ctx, srv.Client(), option.WithEndpoint(srv.URL))
	}
	return paths, scopes
}

func run(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	err := cmd.ExecuteContext(context.Background())
	return out.String(), err
}

func writeBody(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := w.Write([]byte(body)); err != nil {
		t.Errorf("write response: %v", err)
	}
}

func tempFile(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte("payload"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestUploadPrintsFileJSON(t *testing.T) {
	paths, scopes := serve(t, func(w http.ResponseWriter, _ *http.Request) {
		writeBody(t, w, `{"id":"new1","name":"review","mimeType":"application/vnd.google-apps.document","webViewLink":"https://example.org/d","modifiedTime":"2026-01-02T03:04:05Z"}`)
	})
	out, err := run(t, uploadCmd(), tempFile(t, "review.docx"), "--convert", "doc")
	if err != nil {
		t.Fatal(err)
	}
	var f map[string]any
	if err := json.Unmarshal([]byte(out), &f); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}
	if f["id"] != "new1" {
		t.Errorf("output = %s", out)
	}
	if len(*paths) != 1 || (*paths)[0] != "POST /upload/drive/v3/files" {
		t.Errorf("requests = %v", *paths)
	}
	if len(*scopes) != 1 || !strings.HasSuffix((*scopes)[0], "/drive.file") {
		t.Errorf("scopes = %v, want only drive.file", *scopes)
	}
}

func TestUpdatePrintsFileJSON(t *testing.T) {
	paths, _ := serve(t, func(w http.ResponseWriter, _ *http.Request) {
		writeBody(t, w, `{"id":"f1","name":"review","mimeType":"x","modifiedTime":"2026-01-02T03:04:05Z"}`)
	})
	out, err := run(t, updateCmd(), "f1", tempFile(t, "review.docx"), "--convert", "doc")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"id": "f1"`) {
		t.Errorf("output = %s", out)
	}
	if len(*paths) != 1 || (*paths)[0] != "PATCH /upload/drive/v3/files/f1" {
		t.Errorf("requests = %v", *paths)
	}
}

func TestCommentsPrintsArrayEvenWhenEmpty(t *testing.T) {
	_, scopes := serve(t, func(w http.ResponseWriter, _ *http.Request) { writeBody(t, w, `{"comments":[]}`) })
	out, err := run(t, commentsCmd(), "f1")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "[]" {
		t.Errorf("output = %q, want []", out)
	}
	if len(*scopes) != 1 || !strings.HasSuffix((*scopes)[0], "/drive.readonly") {
		t.Errorf("scopes = %v, want only drive.readonly", *scopes)
	}
}

func TestCommentsPrintsComments(t *testing.T) {
	serve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("includeDeleted") != "true" {
			t.Errorf("query = %s, want includeDeleted=true", r.URL.RawQuery)
		}
		writeBody(t, w, `{"comments":[{"id":"c1","content":"Tighten this","quotedFileContent":{"value":"the plan"},"resolved":true,"author":{"displayName":"Reviewer"},"createdTime":"2026-01-02T03:04:05Z","modifiedTime":"2026-01-02T03:04:05Z"}]}`)
	})
	out, err := run(t, commentsCmd(), "f1", "--include-deleted")
	if err != nil {
		t.Fatal(err)
	}
	var got []map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil || len(got) != 1 {
		t.Fatalf("output: %v\n%s", err, out)
	}
	if got[0]["quoted_text"] != "the plan" || got[0]["resolved"] != true {
		t.Errorf("comment = %v", got[0])
	}
}

func TestUsageErrors(t *testing.T) {
	serve(t, func(http.ResponseWriter, *http.Request) { t.Error("usage errors must not reach the API") })
	file := tempFile(t, "r.docx")
	for name, c := range map[string]struct {
		cmd  *cobra.Command
		args []string
	}{
		"upload without a path":    {uploadCmd(), nil},
		"upload with two paths":    {uploadCmd(), []string{file, file}},
		"upload bad convert":       {uploadCmd(), []string{file, "--convert", "pdf"}},
		"update with one argument": {updateCmd(), []string{"f1"}},
		"update bad convert":       {updateCmd(), []string{"f1", file, "--convert", "pdf"}},
		"comments without an ID":   {commentsCmd(), nil},
	} {
		_, err := run(t, c.cmd, c.args...)
		if err == nil || !cmdutil.IsUsage(err) {
			t.Errorf("%s: err = %v, want a usage error", name, err)
		}
	}
}

func TestAPIFailureIsNotAUsageError(t *testing.T) {
	serve(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		writeBody(t, w, `{"error":{"code":404,"message":"File not found"}}`)
	})
	_, err := run(t, commentsCmd(), "missing")
	if err == nil || cmdutil.IsUsage(err) || !strings.Contains(err.Error(), "missing") {
		t.Errorf("err = %v, want a non-usage error naming the file", err)
	}
}
