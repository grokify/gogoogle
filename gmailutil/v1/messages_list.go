package gmailutil

import (
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

const (
	GmailDateFormat  = "2006/01/02"
	ReferenceURL     = "https://developers.google.com/gmail/api/v1/reference"
	TutorialURLGO    = "https://developers.google.com/gmail/api/quickstart/go"
	ListApiReference = "https://developers.google.com/gmail/api/v1/reference/users/messages/list"
	ListApiExample   = "https://stackoverflow.com/questions/43057478/google-api-go-client-listing-messages-w-label-and-fetching-header-fields"
	FilteringExample = "https://developers.google.com/gmail/api/guides/filtering"
	FilterRules      = "https://support.google.com/mail/answer/7190"

	CategoryForums     = "forums"
	CategoryPromotions = "promotions"
	CategorySocial     = "social"
	CategoryUpdates    = "updates"

	ExampleRule1 = "category:promotions older_than:2y"
	ExampleRule2 = "category:updates older_than:2y"
	ExampleRule3 = "category:social older_than:2y"
)

// ?q=in:sent after:1388552400 before:1391230800
// ?q=in:sent after:2014/01/01 before:2014/02/01 ?q=in:sent after:2014/01/01 before:2014/02/01
// "from:someuser@example.com rfc822msgid:<somemsgid@example.com> is:unread"

/*

Warning: All dates used in the search query are interpretted as midnight on that date in the PST timezone. To specify accurate dates for other timezones pass the value in seconds instead:

*/

type MessagesListOpts struct {
	UserId               string
	IncludeSpamTrash     bool
	LabelIds             []string
	MaxResults           uint64
	PageToken            string
	Query                MessagesListQueryOpts
	Fields               []googleapi.Field
	IfNoneMatchEntityTag string
}

func (opts *MessagesListOpts) Condense() {
	opts.UserId = strings.TrimSpace(opts.UserId)
	opts.LabelIds = stringsutil.SliceCondenseSpace(opts.LabelIds, true, false)
	opts.PageToken = strings.TrimSpace(opts.PageToken)
}

func (opts *MessagesListOpts) Inflate() {
	opts.Condense()
	if len(opts.UserId) == 0 {
		opts.UserId = "me"
	}
}

type MessagesListQueryOpts struct {
	Category    string
	In          string
	From        string
	RFC822msgid string
	After       time.Time
	Before      time.Time
	OlderThan   string // #(mdy)
	NewerThan   string // #(mdy)
	Interval    timeutil.Interval
}

func (opts *MessagesListQueryOpts) TrimSpace() {
	opts.Category = strings.TrimSpace(opts.Category)
	opts.From = strings.TrimSpace(opts.From)
	opts.In = strings.TrimSpace(opts.In)
	opts.RFC822msgid = strings.TrimSpace(opts.RFC822msgid)
	opts.OlderThan = strings.TrimSpace(opts.OlderThan)
	opts.NewerThan = strings.TrimSpace(opts.NewerThan)
}

func (opts *MessagesListQueryOpts) Encode() string {
	opts.TrimSpace()
	parts := []string{}

	if len(opts.Category) > 0 {
		parts = append(parts, "category:"+opts.Category)
	}
	if len(opts.From) > 0 {
		parts = append(parts, "from:"+opts.From)
	}
	if len(opts.In) > 0 {
		parts = append(parts, "in:"+opts.In)
	}
	if len(opts.RFC822msgid) > 0 {
		parts = append(parts, "rfc822msgid:"+opts.RFC822msgid)
	}
	if len(opts.OlderThan) > 0 {
		parts = append(parts, "older_than:"+opts.OlderThan)
	}
	if len(opts.NewerThan) > 0 {
		parts = append(parts, "newer_than:"+opts.NewerThan)
	}
	if !timeutil.TimeIsZeroAny(opts.After) {
		parts = append(parts, "after:"+opts.After.Format(GmailDateFormat))
	}
	if !timeutil.TimeIsZeroAny(opts.Before) {
		parts = append(parts, "before:"+opts.Before.Format(GmailDateFormat))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func GetMessagesList(gs *GmailService, opts MessagesListOpts) (*gmail.ListMessagesResponse, error) {
	opts.Inflate()

	userMessagesListCall := gs.UsersService.Messages.List(opts.UserId)
	userMessagesListCall.IncludeSpamTrash(opts.IncludeSpamTrash)
	if len(opts.LabelIds) > 0 {
		userMessagesListCall.LabelIds(opts.LabelIds...)
	}
	if opts.MaxResults > 0 {
		userMessagesListCall.MaxResults(int64(opts.MaxResults))
	}
	if len(opts.PageToken) > 0 {
		userMessagesListCall.PageToken(opts.PageToken)
	}
	q := opts.Query.Encode()
	if len(q) > 0 {
		userMessagesListCall.Q(q)
	}
	if len(opts.Fields) > 0 {
		userMessagesListCall.Fields(opts.Fields...)
	}
	return userMessagesListCall.Do(gs.APICallOptions...)
}

func GetMessagesFrom(gs *GmailService, rfc822 string) (*gmail.ListMessagesResponse, error) {
	opts := MessagesListOpts{
		Query: MessagesListQueryOpts{
			From: rfc822},
	}

	return GetMessagesList(gs, opts)
}

func InflateMessages(gs *GmailService, userId string, msgMetas []*gmail.Message) ([]*gmail.Message, error) {
	msgs := []*gmail.Message{}
	for _, msgMeta := range msgMetas {
		msg, err := GetMessage(gs, userId, msgMeta.Id)
		if err != nil {
			return msgs, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
