// Formatting text with the Google Slides API
// Video: https://www.youtube.com/watch?v=_O2aUCJyCoQ
package slidesutilexamples

import (
	"os"

	su "github.com/grokify/googleutil/slidesutil/v1"
	"github.com/grokify/gotilla/config"
	omg "github.com/grokify/oauth2more/google"
	"github.com/jessevdk/go-flags"
	//"github.com/google/google-api-go-client/slides/v1"
)

type Options struct {
	EnvFile     string `short:"e" long:"env" description:"Env filepath"`
	NewTokenRaw []bool `short:"n" long:"newtoken" description:"Retrieve new token"`
}

func (opt *Options) NewToken() bool {
	if len(opt.NewTokenRaw) > 0 {
		return true
	}
	return false
}

func Setup() (*su.GoogleSlidesService, error) {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		return nil, err
	}

	err = config.LoadDotEnvFirst(opts.EnvFile, os.Getenv("ENV_PATH"))
	if err != nil {
		return nil, err
	}

	googleClient, err := omg.NewClientFileStoreWithDefaults(
		[]byte(os.Getenv(omg.EnvGoogleAppCredentials)),
		[]string{omg.ScopeDrive, omg.ScopePresentations},
		opts.NewToken())
	if err != nil {
		return nil, err
	}

	return su.NewGoogleSlidesService(googleClient)
}
