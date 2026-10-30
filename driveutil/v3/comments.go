package driveutil

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/drive/v3"
)

const (
	commentFields = "nextPageToken,comments(id,content,quotedFileContent(value),anchor," +
		"author(displayName,emailAddress),createdTime,modifiedTime,resolved,deleted," +
		"replies(id,content,action,author(displayName,emailAddress),createdTime))"
	commentPageSize = 100
)

// Author identifies a comment or reply author. EmailAddress is populated
// only when the API returns it.
type Author struct {
	DisplayName  string `json:"display_name"`
	EmailAddress string `json:"email_address,omitempty"`
}

// Reply is a reply to a comment. Action is "", "resolve", or "reopen".
type Reply struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	Action      string    `json:"action,omitempty"`
	Author      Author    `json:"author"`
	CreatedTime time.Time `json:"created_time"`
}

// Comment is a Drive comment with its replies.
type Comment struct {
	ID           string    `json:"id"`
	Content      string    `json:"content"`
	QuotedText   string    `json:"quoted_text,omitempty"`
	Anchor       string    `json:"anchor,omitempty"`
	Author       Author    `json:"author"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
	Resolved     bool      `json:"resolved"`
	Deleted      bool      `json:"deleted"`
	Replies      []Reply   `json:"replies,omitempty"`
}

// ListCommentsOptions filters ListComments.
type ListCommentsOptions struct {
	IncludeDeleted    bool      `json:"include_deleted"`
	StartModifiedTime time.Time `json:"start_modified_time"`
}

// ListComments returns all comments on a file in API order. It follows
// nextPageToken to the end and returns either every page or an error.
func (s *Service) ListComments(ctx context.Context, fileID string, o ListCommentsOptions) ([]Comment, error) {
	var out []Comment
	token := ""
	for {
		call := s.DriveService.Comments.List(fileID).
			Fields(commentFields).
			PageSize(commentPageSize).
			IncludeDeleted(o.IncludeDeleted).
			Context(ctx)
		if !o.StartModifiedTime.IsZero() {
			call = call.StartModifiedTime(o.StartModifiedTime.UTC().Format(time.RFC3339))
		}
		if token != "" {
			call = call.PageToken(token)
		}
		resp, err := call.Do()
		if err != nil {
			return nil, wrapErr("list comments", fileID, err)
		}
		for _, c := range resp.Comments {
			conv, err := fromDriveComment(c)
			if err != nil {
				return nil, fmt.Errorf("driveutil: list comments %s: %w", fileID, err)
			}
			out = append(out, conv)
		}
		if resp.NextPageToken == "" {
			return out, nil
		}
		token = resp.NextPageToken
	}
}

func fromAuthor(u *drive.User) Author {
	if u == nil {
		return Author{}
	}
	return Author{DisplayName: u.DisplayName, EmailAddress: u.EmailAddress}
}

func fromDriveComment(c *drive.Comment) (Comment, error) {
	created, err := parseTime(c.CreatedTime)
	if err != nil {
		return Comment{}, fmt.Errorf("comment %s: createdTime: %w", c.Id, err)
	}
	modified, err := parseTime(c.ModifiedTime)
	if err != nil {
		return Comment{}, fmt.Errorf("comment %s: modifiedTime: %w", c.Id, err)
	}
	out := Comment{
		ID:           c.Id,
		Content:      c.Content,
		Anchor:       c.Anchor,
		Author:       fromAuthor(c.Author),
		CreatedTime:  created,
		ModifiedTime: modified,
		Resolved:     c.Resolved,
		Deleted:      c.Deleted,
	}
	if c.QuotedFileContent != nil {
		out.QuotedText = c.QuotedFileContent.Value
	}
	for _, r := range c.Replies {
		rt, err := parseTime(r.CreatedTime)
		if err != nil {
			return Comment{}, fmt.Errorf("comment %s: reply %s: createdTime: %w", c.Id, r.Id, err)
		}
		out.Replies = append(out.Replies, Reply{
			ID: r.Id, Content: r.Content, Action: r.Action,
			Author: fromAuthor(r.Author), CreatedTime: rt,
		})
	}
	return out, nil
}
