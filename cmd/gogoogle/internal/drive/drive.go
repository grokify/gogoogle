// Package drive provides the drive group command for the gogoogle CLI.
package drive

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/api/drive/v3"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/cmdutil"
	"github.com/grokify/gogoogle/cmd/gogoogle/internal/config"
	driveutil "github.com/grokify/gogoogle/driveutil/v3"
)

// Cmd is the drive group command.
var Cmd = &cobra.Command{
	Use:   "drive",
	Short: "Google Drive utilities",
	Long:  `Commands for uploading files to Google Drive and reading their comments.`,
}

// open creates the Drive service for the given scopes. Tests replace it.
var open = func(ctx context.Context, scopes ...string) (*driveutil.Service, error) {
	c, err := config.NewHTTPClient(ctx, scopes)
	if err != nil {
		return nil, fmt.Errorf("create authenticated client: %w", err)
	}
	return driveutil.NewService(ctx, c)
}

func init() {
	Cmd.AddCommand(uploadCmd(), updateCmd(), commentsCmd())
}

// parseConvert validates the --convert flag value.
func parseConvert(s string) (driveutil.ConvertTo, error) {
	switch c := driveutil.ConvertTo(s); c {
	case driveutil.ConvertNone, driveutil.ConvertDoc, driveutil.ConvertSheet, driveutil.ConvertSlides:
		return c, nil
	}
	return "", cmdutil.Usagef("invalid --convert %q: want doc, sheet, or slides", s)
}

const (
	scopeWrite = drive.DriveFileScope
	scopeRead  = drive.DriveReadonlyScope
)

func uploadCmd() *cobra.Command {
	var name, folder, convert string
	cmd := &cobra.Command{
		Use:   "upload <path>",
		Short: "Upload a local file, optionally converting it to a Google file",
		Long: `Uploads a local file to Drive and prints the created file as JSON.

With --convert doc|sheet|slides the file becomes a Google Doc, Sheet, or
Slides file (.docx/.odt/.txt/.html, .csv/.xlsx/.ods, .pptx/.odp). Without it
the file is stored unchanged.

Needs the drive.file scope.

Example:
  gogoogle drive upload review.docx --convert doc --folder 1AbC`,
		Args: cmdutil.Args(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			to, err := parseConvert(convert)
			if err != nil {
				return err
			}
			svc, err := open(cmd.Context(), scopeWrite)
			if err != nil {
				return err
			}
			f, err := svc.Upload(cmd.Context(), driveutil.UploadRequest{
				Path: args[0], Name: name, FolderID: folder, ConvertTo: to,
			})
			if err != nil {
				return err
			}
			return cmdutil.WriteJSON(cmd.OutOrStdout(), f)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Name of the Drive file (default: the file name, without extension when converting)")
	cmd.Flags().StringVar(&folder, "folder", "", "ID of the parent folder")
	cmd.Flags().StringVar(&convert, "convert", "", "Convert to a Google file: doc, sheet, or slides")
	return cmd
}

func updateCmd() *cobra.Command {
	var convert string
	cmd := &cobra.Command{
		Use:   "update <fileID> <path>",
		Short: "Replace the content of a Drive file, keeping its ID",
		Long: `Replaces the content of an existing Drive file with a local file and prints
the updated file as JSON. The file ID, link, comments, and sharing are kept.

Needs the drive.file scope, which covers files this app created or opened.

Example:
  gogoogle drive update 1AbC review.docx --convert doc`,
		Args: cmdutil.Args(cobra.ExactArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {
			to, err := parseConvert(convert)
			if err != nil {
				return err
			}
			svc, err := open(cmd.Context(), scopeWrite)
			if err != nil {
				return err
			}
			f, err := svc.Update(cmd.Context(), args[0], args[1], to)
			if err != nil {
				return err
			}
			return cmdutil.WriteJSON(cmd.OutOrStdout(), f)
		},
	}
	cmd.Flags().StringVar(&convert, "convert", "", "Convert to a Google file: doc, sheet, or slides")
	return cmd
}

func commentsCmd() *cobra.Command {
	var includeDeleted bool
	cmd := &cobra.Command{
		Use:   "comments <fileID>",
		Short: "List the comments on a Drive file",
		Long: `Lists every comment on a Drive file, with the quoted text it is anchored to,
its resolved state, and its replies, as a JSON array.

Author email addresses appear only when the API returns them; treat them as
personal data. Needs the drive.readonly scope.

Example:
  gogoogle drive comments 1AbC --include-deleted`,
		Args: cmdutil.Args(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := open(cmd.Context(), scopeRead)
			if err != nil {
				return err
			}
			comments, err := svc.ListComments(cmd.Context(), args[0], driveutil.ListCommentsOptions{IncludeDeleted: includeDeleted})
			if err != nil {
				return err
			}
			if comments == nil {
				comments = []driveutil.Comment{}
			}
			return cmdutil.WriteJSON(cmd.OutOrStdout(), comments)
		},
	}
	cmd.Flags().BoolVar(&includeDeleted, "include-deleted", false, "Include deleted comments")
	return cmd
}
