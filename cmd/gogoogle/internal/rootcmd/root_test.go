package rootcmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/cmdutil"
)

func TestExitCode(t *testing.T) {
	for name, tc := range map[string]struct {
		err  error
		want int
	}{
		"success":         {nil, 0},
		"usage":           {cmdutil.Usagef("bad flag"), 2},
		"wrapped usage":   {fmt.Errorf("outer: %w", cmdutil.Usagef("bad flag")), 2},
		"api or io error": {errors.New("boom"), 1},
	} {
		if got := ExitCode(tc.err); got != tc.want {
			t.Errorf("%s: ExitCode = %d, want %d", name, got, tc.want)
		}
	}
}

func TestRegisteredCommands(t *testing.T) {
	for _, path := range [][]string{
		{"drive", "upload"}, {"drive", "update"}, {"drive", "comments"},
		{"docs", "suggestions"}, {"sheets", "records"},
	} {
		if cmd, _, err := rootCmd.Find(path); err != nil || cmd.Name() != path[1] {
			t.Errorf("gogoogle %v not registered: %v", path, err)
		}
	}
}

// A bad flag on a real subcommand is a usage error, exit status 2.
func TestFlagErrorIsUsage(t *testing.T) {
	rootCmd.SetArgs([]string{"drive", "upload", "--nope"})
	rootCmd.SilenceErrors, rootCmd.SilenceUsage = true, true
	t.Cleanup(func() { rootCmd.SilenceErrors, rootCmd.SilenceUsage = false, false })
	if err := rootCmd.Execute(); ExitCode(err) != 2 {
		t.Errorf("err = %v, exit %d, want 2", err, ExitCode(err))
	}
}
