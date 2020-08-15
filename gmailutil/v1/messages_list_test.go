package gmailutil

import (
	"encoding/json"
	"testing"
	"time"
)

var listQueryStringTests = []struct {
	qString  string
	qAfter   string
	qBefore  string
	qOptions MessageListQueryOpts
}{
	{"from:foo@example.com", "", "", MessageListQueryOpts{From: "foo@example.com"}},
	{"in:Inbox", "", "", MessageListQueryOpts{In: "Inbox"}},
	{"after:2016/01/02 before:2019/11/12", "2016-01-02T00:00:00Z", "2019-11-12T00:00:00Z", MessageListQueryOpts{}},
}

// TestGenerateMessageListQueryString creates a gmail query string.
func TestGenerateMessageListQueryString(t *testing.T) {
	for _, tt := range listQueryStringTests {
		qOptions := tt.qOptions
		jsonOpts, err := json.Marshal(qOptions)
		if len(tt.qAfter) > 0 {
			dt, err := time.Parse(time.RFC3339, tt.qAfter)
			if err != nil {
				t.Errorf("gmailutil.GenerateMessageListQueryString('%s') Error: [%v]", jsonOpts, err.Error())
			}
			qOptions.After = dt
		}
		if len(tt.qBefore) > 0 {
			dt, err := time.Parse(time.RFC3339, tt.qBefore)
			if err != nil {
				t.Errorf("gmailutil.GenerateMessageListQueryString('%s') Error: [%v]", jsonOpts, err.Error())
			}
			qOptions.Before = dt
		}
		gotString := GenerateMessageListQueryString(qOptions)
		if gotString != tt.qString {
			if err != nil {
				t.Errorf("gmailutil.GenerateMessageListQueryString() Error: [%v]", err.Error())
			}
			t.Errorf("gmailutil.GenerateMessageListQueryString('%s') Error: want [%v] got [%v]", jsonOpts, tt.qString, gotString)
		}
	}
}
