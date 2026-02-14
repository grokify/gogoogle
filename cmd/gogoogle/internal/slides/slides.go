// Package slides provides the slides group command for the gogoogle CLI.
package slides

import (
	"github.com/spf13/cobra"
)

// Cmd is the slides group command.
var Cmd = &cobra.Command{
	Use:   "slides",
	Short: "Google Slides utilities",
	Long:  `Commands for working with Google Slides presentations.`,
}

func init() {
	Cmd.AddCommand(contentCmd)
}
