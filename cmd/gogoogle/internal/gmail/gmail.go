// Package gmail provides the gmail group command for the gogoogle CLI.
package gmail

import (
	"github.com/spf13/cobra"
)

// Cmd is the gmail group command.
var Cmd = &cobra.Command{
	Use:   "gmail",
	Short: "Gmail utilities",
	Long:  `Commands for working with Gmail.`,
}

func init() {
	Cmd.AddCommand(mergeCmd)
	Cmd.AddCommand(sendMarkdownCmd)
}
