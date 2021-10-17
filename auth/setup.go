// Formatting text with the Google Slides API
// Video: https://www.youtube.com/watch?v=_O2aUCJyCoQ
package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/google"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/jessevdk/go-flags"
	"google.golang.org/api/slides/v1"
	//"github.com/google/google-api-go-client/slides/v1"
)

type Options struct {
	EnvFile     string `short:"e" long:"env" description:"Env filepath"`
	NewTokenRaw []bool `short:"n" long:"newtoken" description:"Retrieve new token"`
}

const (
	apiErrorTokenExpired = "oauth2: token expired and refresh token is not set"
	apiCliTokenRefresh   = ": use `-n` option to refresh token."
)

func (opt *Options) NewToken() bool {
	if len(opt.NewTokenRaw) > 0 {
		return true
	}
	return false
}

func Setup() (*http.Client, error) {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		return nil, err
	}

	err = config.LoadDotEnvFirst(opts.EnvFile, os.Getenv("ENV_PATH"))
	if err != nil {
		return nil, err
	}

	return google.NewClientFileStoreWithDefaults(
		[]byte(os.Getenv(google.EnvGoogleAppCredentials)),
		stringsutil.SplitCondenseSpace(os.Getenv(google.EnvGoogleAppScopes), ","),
		opts.NewToken())
}

func NewGoogleHTTPClient(forceNewToken bool) (*http.Client, error) {
	conf, err := google.ConfigFromEnv(google.ClientSecretEnv,
		[]string{slides.DriveScope, slides.PresentationsScope})
	if err != nil {
		return nil, err
	}

	tokenFile := "slides.googleapis.com-go-quickstart.json"
	tokenStore, err := goauth.NewTokenStoreFileDefault(tokenFile, true, 0700)
	if err != nil {
		return nil, err
	}

	client, err := goauth.NewClientWebTokenStore(
		context.Background(), conf,
		tokenStore, forceNewToken, "mystate")
	if err != nil {
		panic(err)
		return nil, err
	}
	return client, err
}

// WrapError adds CLI instructions to refresh the auth token.
func WrapError(err error) error {
	if strings.Index(err.Error(), apiErrorTokenExpired) > -1 {
		err = fmt.Errorf("%s%s", err.Error(), apiCliTokenRefresh)
	}
	return err
}
