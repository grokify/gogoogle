package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/grokify/gogoogle/gmailutil/v1/mailmerge"
)

func main() {
	if cnt, err := mailmerge.ExecMailMergeCLI(context.Background()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	} else {
		slog.Info(fmt.Sprintf("Successfully sent (%d) email messages", cnt))
		os.Exit(0)
	}
}
