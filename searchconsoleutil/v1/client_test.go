package searchconsoleutil

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"google.golang.org/api/googleapi"
)

// fake is a stand-in for the Search Console API that records requests.
type fake struct {
	srv *httptest.Server

	mu     sync.Mutex
	paths  []string
	bodies []string
	auth   []string
}

// newFake serves routes, keyed by "METHOD path" with the path unescaped.
func newFake(t *testing.T, routes map[string]func(w http.ResponseWriter, body []byte)) *fake {
	t.Helper()
	f := &fake{}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		f.mu.Lock()
		f.paths = append(f.paths, r.Method+" "+r.URL.Path)
		f.bodies = append(f.bodies, string(b))
		f.auth = append(f.auth, r.Header.Get("Authorization"))
		f.mu.Unlock()
		h, ok := routes[r.Method+" "+r.URL.Path]
		if !ok {
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusNotImplemented)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		h(w, b)
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fake) client(t *testing.T) *Client {
	t.Helper()
	c, err := New(context.Background(), f.srv.Client(), WithEndpoint(f.srv.URL+"/"))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func reply(s string) func(http.ResponseWriter, []byte) {
	return func(w http.ResponseWriter, _ []byte) {
		if _, err := io.WriteString(w, s); err != nil {
			panic(err)
		}
	}
}

func fail(code int, reason string) func(http.ResponseWriter, []byte) {
	return func(w http.ResponseWriter, _ []byte) {
		w.WriteHeader(code)
		body := `{"error":{"code":` + itoa(code) + `,"message":"nope","errors":[{"reason":"` + reason + `","message":"nope"}]}}`
		if _, err := io.WriteString(w, body); err != nil {
			panic(err)
		}
	}
}

func itoa(n int) string {
	b, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestNewRejectsNilClient(t *testing.T) {
	if _, err := New(context.Background(), nil); !errors.Is(err, ErrHTTPClientCannotBeNil) {
		t.Fatalf("err = %v, want ErrHTTPClientCannotBeNil", err)
	}
}

func TestErrorsAreClassified(t *testing.T) {
	cases := []struct {
		name   string
		code   int
		reason string
		want   error
	}{
		{"credentials", http.StatusUnauthorized, "authError", ErrUnauthorized},
		{"property", http.StatusForbidden, "forbidden", ErrForbidden},
		{"quota in a 403", http.StatusForbidden, "dailyLimitExceeded", ErrQuotaExceeded},
		{"rate limit", http.StatusTooManyRequests, "rateLimitExceeded", ErrQuotaExceeded},
		{"missing", http.StatusNotFound, "notFound", ErrNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFake(t, map[string]func(http.ResponseWriter, []byte){
				"GET /webmasters/v3/sites": fail(tc.code, tc.reason),
			})
			_, err := f.client(t).Sites(context.Background())
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			var ge *googleapi.Error
			if !errors.As(err, &ge) || ge.Code != tc.code {
				t.Fatalf("underlying error not kept: %v", err)
			}
		})
	}
}

func TestOtherErrorsAreNotClassified(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		"GET /webmasters/v3/sites": fail(http.StatusInternalServerError, "backendError"),
	})
	_, err := f.client(t).Sites(context.Background())
	if err == nil {
		t.Fatal("want an error")
	}
	for _, k := range []error{ErrUnauthorized, ErrForbidden, ErrNotFound, ErrQuotaExceeded} {
		if errors.Is(err, k) {
			t.Fatalf("a 500 was classified as %v", k)
		}
	}
}

func TestCancelledContextIsNotWrapped(t *testing.T) {
	f := newFake(t, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := f.client(t).Sites(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}

// testKey returns a service account key whose token URI is tokenURI.
func testKey(t *testing.T, tokenURI string) []byte {
	t.Helper()
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		t.Fatal(err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	b, err := json.Marshal(map[string]string{
		"type":           "service_account",
		"client_email":   "reader@example-project.iam.gserviceaccount.com",
		"private_key":    string(pemKey),
		"private_key_id": "key-1",
		"token_uri":      tokenURI,
	})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestServiceAccountEmail(t *testing.T) {
	got, err := ServiceAccountEmail(testKey(t, "https://example.test/token"))
	if err != nil || got != "reader@example-project.iam.gserviceaccount.com" {
		t.Fatalf("got %q, %v", got, err)
	}
	for name, in := range map[string]string{
		"not json":   "nope",
		"other type": `{"type":"authorized_user","client_email":"x@y"}`,
		"no email":   `{"type":"service_account"}`,
		"empty":      ``,
	} {
		if _, err := ServiceAccountEmail([]byte(in)); !errors.Is(err, ErrNotServiceAccountKey) {
			t.Errorf("%s: err = %v, want ErrNotServiceAccountKey", name, err)
		}
	}
}

func TestNewFromServiceAccountKeyRejectsOtherCredentials(t *testing.T) {
	_, err := NewFromServiceAccountKey(context.Background(), []byte(`{"type":"authorized_user"}`))
	if !errors.Is(err, ErrNotServiceAccountKey) {
		t.Fatalf("err = %v, want ErrNotServiceAccountKey", err)
	}
}

// TestServiceAccountSignsInWithReadOnlyScope runs the whole sign-in against
// a fake token endpoint: the key is exchanged for a token that carries only
// the read-only scope, and the API call carries that token.
func TestServiceAccountSignsInWithReadOnlyScope(t *testing.T) {
	var tokenForm string
	var mu sync.Mutex
	tokens := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read token request: %v", err)
		}
		mu.Lock()
		tokenForm = string(b)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, `{"access_token":"test-access-token","token_type":"Bearer","expires_in":3600}`); err != nil {
			t.Errorf("write token: %v", err)
		}
	}))
	t.Cleanup(tokens.Close)

	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		"GET /webmasters/v3/sites": reply(`{"siteEntry":[{"siteUrl":"sc-domain:example.com","permissionLevel":"siteFullUser"}]}`),
	})
	c, err := NewFromServiceAccountKey(context.Background(), testKey(t, tokens.URL), WithEndpoint(f.srv.URL+"/"))
	if err != nil {
		t.Fatal(err)
	}
	sites, err := c.Sites(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(sites) != 1 || sites[0].URL != "sc-domain:example.com" {
		t.Fatalf("sites = %+v", sites)
	}
	if got := f.auth[0]; got != "Bearer test-access-token" {
		t.Errorf("Authorization = %q", got)
	}
	mu.Lock()
	defer mu.Unlock()
	if tokenForm == "" || !strings.Contains(tokenForm, "grant_type=") {
		t.Fatalf("no token request was made: %q", tokenForm)
	}
}
