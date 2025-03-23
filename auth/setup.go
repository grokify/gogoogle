// Formatting text with the Google Slides API
// Video: https://www.youtube.com/watch?v=_O2aUCJyCoQ
package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/google"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/type/stringsutil"
	flags "github.com/jessevdk/go-flags"
	slides "google.golang.org/api/slides/v1"
	//"github.com/google/google-api-go-client/slides/v1"
)

type Options struct {
	EnvFile     string `short:"e" long:"env" description:"Env filepath"`
	NewTokenRaw []bool `short:"n" long:"newtoken" description:"Retrieve new token"`
}

const (
	apiErrorTokenExpired = "oauth2: token expired and refresh token is not set" // #nosec G101
	apiCliTokenRefresh   = ": use `-n` option to refresh token."                // #nosec G101
)

func (opt *Options) NewToken() bool {
	return len(opt.NewTokenRaw) > 0
}

func Setup(ctx context.Context) (*http.Client, error) {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		return nil, err
	}

	_, err = config.LoadDotEnv([]string{opts.EnvFile, os.Getenv("ENV_PATH")}, 1)
	if err != nil {
		return nil, err
	}

	return google.NewClientFileStoreWithDefaults(
		ctx,
		[]byte(os.Getenv(google.EnvGoogleAppCredentials)),
		stringsutil.SplitTrimSpace(os.Getenv(google.EnvGoogleAppScopes), ",", true),
		opts.NewToken())
}

func NewGoogleHTTPClient(forceNewToken bool) (*http.Client, error) {
	conf, err := google.ConfigFromEnv(google.ClientSecretEnv,
		[]string{slides.DriveScope, slides.PresentationsScope})
	if err != nil {
		return nil, err
	}

	tokenFile := "slides.googleapis.com-go-quickstart.json"
	tokenStore, err := authutil.NewTokenStoreFileDefault(tokenFile, true, 0700)
	if err != nil {
		return nil, err
	}

	client, err := authutil.NewClientWebTokenStore(
		context.Background(), conf,
		tokenStore, forceNewToken, "mystate")
	if err != nil {
		return nil, err
	}
	return client, err
}

// WrapError adds CLI instructions to refresh the auth token.
func WrapError(err error) error {
	if strings.Contains(err.Error(), apiErrorTokenExpired) {
		err = fmt.Errorf("%s%s", err.Error(), apiCliTokenRefresh)
	}
	return err
}
