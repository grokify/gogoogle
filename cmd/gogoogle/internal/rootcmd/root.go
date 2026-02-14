// Package rootcmd provides the root command for the gogoogle CLI.
package rootcmd

import (
	"github.com/spf13/cobra"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/config"
	"github.com/grokify/gogoogle/cmd/gogoogle/internal/gmail"
	"github.com/grokify/gogoogle/cmd/gogoogle/internal/slides"
)

var (
	// Persistent flags for authentication.
	credentials           string
	goauthCredentialsFile string
	goauthCredentialsAcct string
)

var rootCmd = &cobra.Command{
	Use:   "gogoogle",
	Short: "Unified CLI for Google API utilities",
	Long: `gogoogle is a unified command-line interface for working with Google APIs.

It provides subcommands for working with Google Slides, Gmail, and other
Google services.

Authentication:
  Use one of the following authentication methods:

  1. Google service account credentials:
     --credentials /path/to/service-account.json

  2. goauth CredentialsSet file:
     --goauth-credentials-file /path/to/credentials.json \
     --goauth-credentials-account myaccount

Environment variables:
  GOOGLE_CREDENTIALS_FILE       - Default for --credentials
  GOAUTH_CREDENTIALS_FILE       - Default for --goauth-credentials-file
  GOAUTH_CREDENTIALS_ACCOUNT    - Default for --goauth-credentials-account`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set credentials in config package for subcommands to access.
		config.SetCredentials(credentials, goauthCredentialsFile, goauthCredentialsAcct)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&credentials, "credentials", "",
		"Path to Google service account credentials JSON file (env: GOOGLE_CREDENTIALS_FILE)")
	rootCmd.PersistentFlags().StringVar(&goauthCredentialsFile, "goauth-credentials-file", "",
		"Path to goauth CredentialsSet JSON file (env: GOAUTH_CREDENTIALS_FILE)")
	rootCmd.PersistentFlags().StringVar(&goauthCredentialsAcct, "goauth-credentials-account", "",
		"Account key within goauth CredentialsSet file (env: GOAUTH_CREDENTIALS_ACCOUNT)")

	// Add subcommands.
	rootCmd.AddCommand(slides.Cmd)
	rootCmd.AddCommand(gmail.Cmd)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
