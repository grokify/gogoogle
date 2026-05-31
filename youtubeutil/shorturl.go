package youtubeutil

import (
	"errors"
	"net/url"
	"strings"

	"github.com/grokify/mogo/net/urlutil"
)

const (
	ShortURLHost = "https://youtu.be/"
)

func ShortURL(vid string) (string, error) {
	vid = strings.TrimSpace(vid)
	if strings.Index(vid, ShortURLHost) == 0 && len(vid) > len(ShortURLHost) {
		return vid, nil
	}
	if urlutil.IsHTTP(vid, true, true) {
		if strings.Contains(vid, "?") {
			if u, err := url.Parse(vid); err == nil {
				qry := u.Query()
				vid := qry.Get("v")
				if len(vid) > 0 {
					return urlutil.JoinAbsolute(ShortURLHost, vid), nil
				} else {
					return vid, errors.New("url has query without `v`")
				}
			} else {
				return vid, err
			}
		} else {
			return vid, errors.New("url has query without `v`")
		}
	}
	return urlutil.JoinAbsolute(ShortURLHost, vid), nil
}
