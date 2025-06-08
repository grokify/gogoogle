package mailmerge

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"

	"github.com/grokify/gocharts/v2/data/table"
	gmailutil "github.com/grokify/gogoogle/gmailutil/v1"
	"github.com/grokify/gogoogle/sheetsutil/iwark"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/mime/multipartutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/mailutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/sogo/text/mustacheutil"
)

const (
	ColumnTo   = "TO"
	ColumnCc   = "CC"
	ColumnBcc  = "BCC"
	ColumnFrom = "FROM" // can be "me"

	templateTypeBodyHTML    = "bodyhtml"
	templateTypeBodyText    = "bodytext"
	templateTypeSubjectText = "subjecttext"
)

var ErrMailMergeOptsCannotBeNil = errors.New("parameter MailMergeOpts cannot be nil")

type MailMergeOpts struct {
	GoauthCredsFile                 string   `short:"c" long:"goauth-credentials-file" description:"The Google Sheet ID" required:"true"`
	GoauthAccountKey                string   `short:"k" long:"goauth-account-key" description:"The Google Sheet ID"`
	RecipientsGoogleSheetID         string   `short:"s" long:"sheet-id" description:"The Google Sheet ID"`
	RecipientsGoogleSheetIndex      uint     `short:"x" long:"sheet-index" description:"The Google Sheet Index"`
	RecipientsGoogleSheetHeaderRows uint32   `short:"r" long:"sheet-header-row-count" description:"The Google Sheet header row count"`
	SubjectTemplateTextFilename     string   `short:"j" long:"subject-template" description:"Subject template"`
	BodyTemplateHTMLFilename        string   `long:"html-template" description:"Body tmeplate for HTML"`
	BodyTemplateTextFilename        string   `short:"t" long:"text-template" description:"Body template for text"`
	InlineFilenames                 []string `short:"i" long:"inline-filename" description:"Inline filenames"`
	AttachmentsFilenames            []string `short:"a" long:"attachment-filename" description:"Filenames as attachments"`

	GoogleClient       *http.Client
	BodyCommonPartsSet multipartutil.PartsSet
}

func (opts MailMergeOpts) Validate() error {
	var errorMsgs []string
	if opts.GoogleClient == nil {
		errorMsgs = append(errorMsgs, "GoogleClient is nil")
	}
	if strings.TrimSpace(opts.RecipientsGoogleSheetID) == "" {
		errorMsgs = append(errorMsgs, "RecipientsGoogleSheetID is empty")
	}
	if opts.RecipientsGoogleSheetHeaderRows == 0 {
		errorMsgs = append(errorMsgs, "RecipientsGoogleSheetHeaderRows cannot be 0")
	}
	if strings.TrimSpace(opts.SubjectTemplateTextFilename) == "" {
		errorMsgs = append(errorMsgs, "subject templates is empty: SubjectTemplateFilename")
	}
	if strings.TrimSpace(opts.BodyTemplateHTMLFilename) == "" && strings.TrimSpace(opts.BodyTemplateTextFilename) == "" {
		errorMsgs = append(errorMsgs, "body templates are both empty: BodyTemplateHTMLFilename and BodyTemplateTextFilename")
	}
	if len(errorMsgs) > 0 {
		return fmt.Errorf("errors: (%s)", strings.Join(errorMsgs, ", "))
	} else {
		return nil
	}
}

type MailMerge struct {
	BodyTemplateSet *mustacheutil.MustacheSet
	Table           *table.Table
	CommonPartsSet  multipartutil.PartsSet
	GmailService    *gmailutil.GmailService
}

func NewMailMerge(ctx context.Context, opts *MailMergeOpts) (*MailMerge, error) {
	if opts == nil {
		return nil, ErrMailMergeOptsCannotBeNil
	} else if err := opts.Validate(); err != nil {
		return nil, err
	}
	mm := MailMerge{
		BodyTemplateSet: &mustacheutil.MustacheSet{
			Filenames: map[string]string{
				templateTypeBodyHTML:    opts.BodyTemplateHTMLFilename,
				templateTypeBodyText:    opts.BodyTemplateTextFilename,
				templateTypeSubjectText: opts.SubjectTemplateTextFilename,
			},
		},
		CommonPartsSet: opts.BodyCommonPartsSet.Clone(),
	}

	if err := mm.BodyTemplateSet.ReadTemplates(); err != nil {
		return nil, err
	} else if err := mm.loadFilesInline(opts.InlineFilenames); err != nil {
		return nil, err
	} else if err := mm.loadFilesAttachment(opts.AttachmentsFilenames); err != nil {
		return nil, err
	}

	if strings.TrimSpace(opts.RecipientsGoogleSheetID) != "" {
		if opts.GoogleClient == nil {
			return nil, errors.New("google client cannot be nil with google sheet id")
		}
		tbl, err := iwark.ReadTableFromClient(
			opts.GoogleClient,
			opts.RecipientsGoogleSheetID,
			opts.RecipientsGoogleSheetIndex,
			opts.RecipientsGoogleSheetHeaderRows)
		if err != nil {
			return nil, err
		} else {
			mm.Table = tbl
		}
	}

	if gmSvc, err := gmailutil.NewGmailService(ctx, opts.GoogleClient); err != nil {
		return nil, err
	} else {
		mm.GmailService = gmSvc
	}

	return &mm, nil
}

func (mm *MailMerge) loadFilesAttachment(filenames []string) error {
	return mm.loadFiles(httputilmore.DispositionTypeAttachment, filenames)
}

func (mm *MailMerge) loadFilesInline(filenames []string) error {
	return mm.loadFiles(httputilmore.DispositionTypeInline, filenames)
}

func (mm *MailMerge) loadFiles(dispositionType string, filenames []string) error {
	seen := map[string]bool{}
	for _, filename := range filenames {
		_, filenameonly := filepath.Split(filename)
		if _, ok := seen[filenameonly]; ok {
			return errors.New("duplicate inline filename")
		} else {
			seen[filenameonly] = true
		}
		mm.CommonPartsSet.Parts = append(mm.CommonPartsSet.Parts, multipartutil.Part{
			Type:             multipartutil.PartTypeFilepath,
			BodyEncodeBase64: true,
			BodyDataFilepath: filenameonly,
			DispositionType:  dispositionType,
			HeaderRaw: textproto.MIMEHeader{
				httputilmore.HeaderContentID: []string{filenameonly},
			},
		})
	}
	return nil
}

func (mm *MailMerge) Messages() ([]mailutil.MessageWriter, error) {
	var msgs []mailutil.MessageWriter
	if mm.BodyTemplateSet == nil {
		return msgs, errors.New("template set cannot be nil")
	} else if mm.Table == nil {
		return msgs, errors.New("recipient table cannot be nil")
	} else if len(mm.Table.Rows) == 0 {
		return msgs, errors.New("table has no recipients")
	}
	tbl := mm.Table

	for i, row := range tbl.Rows {
		rowTry := stringsutil.SliceCondenseSpace(row, true, false)
		if len(rowTry) == 0 {
			continue
		}
		rowMap := tbl.Columns.RowMap(row, false)

		toAddrs, err := mailutil.ParseAddressList(tbl.Columns.MustCellString(ColumnTo, row))
		if err != nil {
			return msgs, err
		}
		ccAddrs, err := mailutil.ParseAddressList(tbl.Columns.MustCellString(ColumnCc, row))
		if err != nil {
			return msgs, err
		}
		bccAddrs, err := mailutil.ParseAddressList(tbl.Columns.MustCellString(ColumnBcc, row))
		if err != nil {
			return msgs, err
		}
		if len(toAddrs.FilterInclWithoutAddress()) > 0 {
			return msgs, fmt.Errorf("to addresses include empty (%s)", tbl.Columns.MustCellString(ColumnTo, row))
		}
		if len(ccAddrs.FilterInclWithoutAddress()) > 0 {
			return msgs, fmt.Errorf("cc addresses include empty (%s)", tbl.Columns.MustCellString(ColumnCc, row))
		}
		if len(bccAddrs.FilterInclWithoutAddress()) > 0 {
			return msgs, fmt.Errorf("bcc addresses include empty (%s)", tbl.Columns.MustCellString(ColumnBcc, row))
		}

		bytesSubject, err := mm.BodyTemplateSet.RenderTemplateOrDefault(templateTypeSubjectText, rowMap, []byte{})
		if err != nil {
			return msgs, err
		}
		bytesBodyText, err := mm.BodyTemplateSet.RenderTemplateOrDefault(templateTypeBodyText, rowMap, []byte{})
		if err != nil {
			return msgs, err
		}
		bytesBodyHTML, err := mm.BodyTemplateSet.RenderTemplateOrDefault(templateTypeBodyHTML, rowMap, []byte{})
		if err != nil {
			return msgs, err
		}

		msgParts, err := multipartutil.NewPartsSetMail(bytesBodyText, bytesBodyHTML, mm.CommonPartsSet.Parts)
		if err != nil {
			return msgs, err
		}

		msgout := mailutil.MessageWriter{
			To:           toAddrs,
			Cc:           ccAddrs,
			Bcc:          bccAddrs,
			Subject:      string(bytesSubject),
			BodyPartsSet: msgParts,
		}
		if msgout.RecipientCount() <= 0 {
			if out, err := jsonutil.MarshalSlice(row); err != nil {
				return msgs, err
			} else {
				return msgs, fmt.Errorf("no recpients on row (%d) with data (%s)", i, string(out))
			}
		} else {
			msgs = append(msgs, msgout)
		}
	}

	return msgs, nil
}

func (mm *MailMerge) Send(ctx context.Context, userID string) (int, error) {
	if userID = strings.TrimSpace(userID); userID == "" {
		userID = gmailutil.UserIDMe
	}

	msgs, err := mm.Messages()
	if err != nil {
		return -1, err
	}

	for _, msg := range msgs {
		if _, err := mm.GmailService.Send(ctx, userID, msg); err != nil {
			return -1, err
		}
	}
	return len(msgs), nil
}
