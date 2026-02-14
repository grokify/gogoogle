// gogoogle is a unified CLI for Google API utilities.
package main

import (
	"os"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/rootcmd"
)

func main() {
	if err := rootcmd.Execute(); err != nil {
		os.Exit(1)
	}
}
