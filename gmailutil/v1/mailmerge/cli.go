package mailmerge

import (
	"context"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
	"github.com/jessevdk/go-flags"
)

func ExecMailMergeCLI(ctx context.Context) (int, error) {
	opts := MailMergeOpts{}
	_, err := flags.Parse(&opts)
	if err != nil {
		return -1, err
	}

	if creds, err := goauth.ReadCredentialsFromSetFile(opts.GoauthCredsFile, opts.GoauthAccountKey, true); err != nil {
		return -1, err
	} else if tok, err := creds.NewOrExistingValidToken(ctx); err != nil {
		return -1, err
	} else {
		opts.GoogleClient = authutil.NewClientTokenOAuth2(tok)
	}

	mm, err := NewMailMerge(ctx, &opts)
	if err != nil {
		return -1, err
	}

	return mm.Send(ctx, "")
}
