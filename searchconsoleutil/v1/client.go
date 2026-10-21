package searchconsoleutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	searchconsole "google.golang.org/api/searchconsole/v1"
)

// ScopeReadonly is the only scope this package requests.
const ScopeReadonly = searchconsole.WebmastersReadonlyScope

// Errors returned by Client methods and constructors. A failed API call wraps
// both one of these and the underlying *googleapi.Error, so callers can test
// for either with errors.Is and errors.As.
var (
	// ErrHTTPClientCannotBeNil is returned by New for a nil client.
	ErrHTTPClientCannotBeNil = errors.New("searchconsole: http client cannot be nil")
	// ErrNotServiceAccountKey is returned when credentials are not a
	// service account key.
	ErrNotServiceAccountKey = errors.New("searchconsole: credentials are not a service account key")
	// ErrUnauthorized means Google rejected the credentials (HTTP 401).
	ErrUnauthorized = errors.New("searchconsole: credentials rejected")
	// ErrForbidden means the credentials cannot access the property
	// (HTTP 403), for example because the service account was not added to it.
	ErrForbidden = errors.New("searchconsole: no access to the property")
	// ErrNotFound means the property, sitemap or URL does not exist for
	// the credentials (HTTP 404).
	ErrNotFound = errors.New("searchconsole: not found")
	// ErrQuotaExceeded means a rate or daily quota was reached
	// (HTTP 429, or 403 with a quota reason).
	ErrQuotaExceeded = errors.New("searchconsole: quota exceeded")
)

// quotaReasons are the reasons Google gives a 403 that is really a quota.
var quotaReasons = map[string]bool{
	"rateLimitExceeded":     true,
	"userRateLimitExceeded": true,
	"quotaExceeded":         true,
	"dailyLimitExceeded":    true,
}

// Client reads Search Console. Its methods are safe for concurrent use.
type Client struct {
	svc *searchconsole.Service
}

// Option configures a Client.
type Option func(*config)

type config struct{ endpoint string }

// WithEndpoint replaces the API's base URL, for tests and proxies.
func WithEndpoint(url string) Option { return func(c *config) { c.endpoint = url } }

// New returns a Client that sends requests through hc, which must be
// authorized for [ScopeReadonly].
func New(ctx context.Context, hc *http.Client, opts ...Option) (*Client, error) {
	if hc == nil {
		return nil, ErrHTTPClientCannotBeNil
	}
	var cfg config
	for _, o := range opts {
		o(&cfg)
	}
	so := []option.ClientOption{option.WithHTTPClient(hc)}
	if cfg.endpoint != "" {
		so = append(so, option.WithEndpoint(cfg.endpoint))
	}
	svc, err := searchconsole.NewService(ctx, so...)
	if err != nil {
		return nil, fmt.Errorf("searchconsole: create service: %w", err)
	}
	return &Client{svc: svc}, nil
}

// NewFromServiceAccountKey returns a Client that authenticates as the
// service account in key, the JSON file Google issues for it. Only
// service account keys are accepted, so a key from an untrusted source
// cannot change what the client authenticates as.
//
// ctx is used whenever the access token is refreshed, so it must outlive the
// Client; pass a long-lived context, not a request's.
func NewFromServiceAccountKey(ctx context.Context, key []byte, opts ...Option) (*Client, error) {
	if _, err := ServiceAccountEmail(key); err != nil {
		return nil, err
	}
	cfg, err := google.JWTConfigFromJSON(key, ScopeReadonly)
	if err != nil {
		return nil, fmt.Errorf("searchconsole: read service account key: %w", err)
	}
	return New(ctx, cfg.Client(ctx), opts...)
}

// ServiceAccountEmail returns the email of the service account in key,
// which is the address to add as a user on a Search Console property. It
// returns [ErrNotServiceAccountKey] for any other kind of credentials.
func ServiceAccountEmail(key []byte) (string, error) {
	var k struct {
		Type        string `json:"type"`
		ClientEmail string `json:"client_email"`
	}
	if err := json.Unmarshal(key, &k); err != nil {
		return "", fmt.Errorf("%w: %w", ErrNotServiceAccountKey, err)
	}
	if k.Type != "service_account" || k.ClientEmail == "" {
		return "", ErrNotServiceAccountKey
	}
	return k.ClientEmail, nil
}

// wrap classifies an API error. Errors that are not API errors, such as a
// cancelled context, are returned as they are.
func wrap(err error) error {
	if err == nil {
		return nil
	}
	var ge *googleapi.Error
	if !errors.As(err, &ge) {
		return err
	}
	var kind error
	switch ge.Code {
	case http.StatusUnauthorized:
		kind = ErrUnauthorized
	case http.StatusNotFound:
		kind = ErrNotFound
	case http.StatusTooManyRequests:
		kind = ErrQuotaExceeded
	case http.StatusForbidden:
		kind = ErrForbidden
		for _, e := range ge.Errors {
			if quotaReasons[e.Reason] {
				kind = ErrQuotaExceeded
				break
			}
		}
	default:
		return err
	}
	return fmt.Errorf("%w: %w", kind, err)
}
