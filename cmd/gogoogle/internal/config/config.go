// Package config provides shared configuration for the gogoogle CLI.
package config

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/google"
)

var (
	// Persistent flags for authentication.
	credentials           string
	goauthCredentialsFile string
	goauthCredentialsAcct string
	mu                    sync.RWMutex
)

// ErrNoCredentials is returned when no authentication credentials are provided.
var ErrNoCredentials = errors.New("credentials required: use --credentials or --goauth-credentials-file with --goauth-credentials-account")

// ErrMultipleCredentials is returned when both credential methods are provided.
var ErrMultipleCredentials = errors.New("cannot use both --credentials and --goauth-credentials-file")

// SetCredentials sets the credential values (called from root command).
func SetCredentials(creds, goauthFile, goauthAcct string) {
	mu.Lock()
	defer mu.Unlock()
	credentials = creds
	goauthCredentialsFile = goauthFile
	goauthCredentialsAcct = goauthAcct
}

// GetCredentials returns the current credential flag values.
func GetCredentials() (creds, goauthFile, goauthAcct string) {
	mu.RLock()
	defer mu.RUnlock()
	return credentials, goauthCredentialsFile, goauthCredentialsAcct
}

// NewHTTPClient creates an authenticated HTTP client for the specified Google API scopes.
// It uses the authentication flags set on the root command.
func NewHTTPClient(ctx context.Context, scopes []string) (*http.Client, error) {
	creds, goauthFile, goauthAcct := GetCredentials()

	// Apply environment variable defaults.
	if creds == "" {
		creds = os.Getenv("GOOGLE_CREDENTIALS_FILE")
	}
	if goauthFile == "" {
		goauthFile = os.Getenv("GOAUTH_CREDENTIALS_FILE")
	}
	if goauthAcct == "" {
		goauthAcct = os.Getenv("GOAUTH_CREDENTIALS_ACCOUNT")
	}

	hasGoogleCreds := creds != ""
	hasGoauthCreds := goauthFile != "" && goauthAcct != ""

	if !hasGoogleCreds && !hasGoauthCreds {
		return nil, ErrNoCredentials
	}

	if hasGoogleCreds && hasGoauthCreds {
		return nil, ErrMultipleCredentials
	}

	if hasGoauthCreds {
		return goauth.NewClient(ctx, goauthFile, goauthAcct)
	}

	return google.NewClientSvcAccountFromFile(ctx, creds, scopes...)
}
