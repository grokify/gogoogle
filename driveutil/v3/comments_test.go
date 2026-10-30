package driveutil

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestListCommentsPaging(t *testing.T) {
	pages := []string{
		`{"nextPageToken":"p2","comments":[{"id":"c1","content":"one","anchor":"a1","quotedFileContent":{"value":"quoted"},
		  "author":{"displayName":"Ann","emailAddress":"ann@example.com"},"createdTime":"2026-01-01T00:00:00Z","modifiedTime":"2026-01-02T00:00:00Z",
		  "replies":[{"id":"r1","content":"done","action":"resolve","author":{"displayName":"Bob"},"createdTime":"2026-01-03T00:00:00Z"}],"resolved":true},
		  {"id":"c2","content":"two"}]}`,
		`{"nextPageToken":"p3","comments":[{"id":"c3","content":"three","deleted":true}]}`,
		`{"comments":[{"id":"c4","content":"four"}]}`,
	}
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, n int) { writeBody(t, w, pages[n-1]) })
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	got, err := f.service(t).ListComments(context.Background(), "file1", ListCommentsOptions{IncludeDeleted: true, StartModifiedTime: start})
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, c := range got {
		ids = append(ids, c.ID)
	}
	if len(ids) != 4 || ids[0] != "c1" || ids[1] != "c2" || ids[2] != "c3" || ids[3] != "c4" {
		t.Fatalf("ids = %v", ids)
	}
	c := got[0]
	if c.QuotedText != "quoted" || c.Anchor != "a1" || !c.Resolved || c.Author.DisplayName != "Ann" || c.Author.EmailAddress != "ann@example.com" {
		t.Fatalf("unexpected comment: %+v", c)
	}
	if c.CreatedTime.Day() != 1 || c.ModifiedTime.Day() != 2 {
		t.Fatalf("bad times: %+v", c)
	}
	if len(c.Replies) != 1 || c.Replies[0].Action != "resolve" || c.Replies[0].Author.DisplayName != "Bob" || c.Replies[0].CreatedTime.Day() != 3 {
		t.Fatalf("unexpected replies: %+v", c.Replies)
	}
	if !got[2].Deleted || got[1].Author != (Author{}) {
		t.Fatalf("unexpected: %+v %+v", got[1], got[2])
	}

	reqs := f.requests()
	if len(reqs) != 3 {
		t.Fatalf("want 3 requests, got %d", len(reqs))
	}
	wantTokens := []string{"", "p2", "p3"}
	for i, r := range reqs {
		if r.Method != "GET" || r.Path != "/files/file1/comments" {
			t.Fatalf("req %d: %s %s", i, r.Method, r.Path)
		}
		q := r.Query
		if q["fields"][0] != commentFields || q["pageSize"][0] != "100" || q["includeDeleted"][0] != "true" ||
			q["startModifiedTime"][0] != "2026-01-01T12:00:00Z" {
			t.Fatalf("req %d: query %v", i, q)
		}
		tok := ""
		if v := q["pageToken"]; len(v) > 0 {
			tok = v[0]
		}
		if tok != wantTokens[i] {
			t.Fatalf("req %d: pageToken %q, want %q", i, tok, wantTokens[i])
		}
	}
}

func TestListCommentsDefaults(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, `{}`) })
	got, err := f.service(t).ListComments(context.Background(), "file1", ListCommentsOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
	q := f.requests()[0].Query
	if _, ok := q["startModifiedTime"]; ok {
		t.Fatalf("startModifiedTime should be absent: %v", q)
	}
	if v := q["includeDeleted"]; len(v) != 0 && v[0] != "false" {
		t.Fatalf("includeDeleted = %v", v)
	}
}

func TestListCommentsPageErrorIsNotPartial(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, n int) {
		if n == 1 {
			writeBody(t, w, `{"nextPageToken":"p2","comments":[{"id":"c1"}]}`)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		writeBody(t, w, `{"error":{"code":403,"message":"no"}}`)
	})
	got, err := f.service(t).ListComments(context.Background(), "file1", ListCommentsOptions{})
	if !errors.Is(err, ErrPermissionDenied) || got != nil {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestListCommentsNotFound(t *testing.T) {
	f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(http.StatusNotFound)
		writeBody(t, w, `{"error":{"code":404,"message":"gone"}}`)
	})
	if _, err := f.service(t).ListComments(context.Background(), "x", ListCommentsOptions{}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestListCommentsBadTimes(t *testing.T) {
	bodies := []string{
		`{"comments":[{"id":"c1","createdTime":"bad"}]}`,
		`{"comments":[{"id":"c1","modifiedTime":"bad"}]}`,
		`{"comments":[{"id":"c1","replies":[{"id":"r1","createdTime":"bad"}]}]}`,
	}
	for _, b := range bodies {
		f := newFake(t, func(w http.ResponseWriter, _ *http.Request, _ int) { writeBody(t, w, b) })
		if _, err := f.service(t).ListComments(context.Background(), "file1", ListCommentsOptions{}); err == nil {
			t.Fatalf("expected error for %s", b)
		}
	}
}
