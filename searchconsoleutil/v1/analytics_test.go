package searchconsoleutil

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

const queryPath = "POST /webmasters/v3/sites/sc-domain:example.com/searchAnalytics/query"

// day is a date in 2026 at midday, to show that the time of day is ignored.
func day(m time.Month, d int) time.Time { return time.Date(2026, m, d, 13, 45, 0, 0, time.UTC) }

func TestQuerySendsTheRequest(t *testing.T) {
	var sent map[string]any
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		queryPath: func(w http.ResponseWriter, body []byte) {
			if err := json.Unmarshal(body, &sent); err != nil {
				t.Errorf("request is not JSON: %v", err)
			}
			if _, err := w.Write([]byte(`{
				"rows":[{"keys":["a","b"],"clicks":3,"impressions":40,"ctr":0.075,"position":4.5}],
				"metadata":{"firstIncompleteDate":"2026-10-19"}}`)); err != nil {
				t.Error(err)
			}
		},
	})
	res, err := f.client(t).Query(context.Background(), "sc-domain:example.com", Query{
		Start: day(9, 22), End: day(10, 19),
		Dimensions: []Dimension{DimensionQuery, DimensionPage},
		SearchType: SearchTypeWeb, DataState: DataStateAll,
		Filters:  []Filter{{Dimension: DimensionCountry, Operator: FilterEquals, Expression: "usa"}},
		RowLimit: 50, StartRow: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if sent["startDate"] != "2026-09-22" || sent["endDate"] != "2026-10-19" {
		t.Errorf("dates = %v..%v", sent["startDate"], sent["endDate"])
	}
	if sent["searchType"] != "web" || sent["dataState"] != "all" || sent["rowLimit"] != float64(50) || sent["startRow"] != float64(100) {
		t.Errorf("request = %v", sent)
	}
	if dims, _ := sent["dimensions"].([]any); len(dims) != 2 || dims[0] != "query" || dims[1] != "page" {
		t.Errorf("dimensions = %v", sent["dimensions"])
	}
	groups, _ := sent["dimensionFilterGroups"].([]any)
	if len(groups) != 1 || !strings.Contains(string(mustJSON(t, groups)), `"expression":"usa"`) {
		t.Errorf("filters = %v", groups)
	}
	if len(res.Rows) != 1 || res.Rows[0].Clicks != 3 || res.Rows[0].Keys[1] != "b" || res.Rows[0].CTR != 0.075 || res.Rows[0].Position != 4.5 {
		t.Errorf("rows = %+v", res.Rows)
	}
	if got := res.FirstIncompleteDate.Format("2006-01-02"); got != "2026-10-19" {
		t.Errorf("first incomplete date = %s", got)
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestQueryValidation(t *testing.T) {
	f := newFake(t, nil)
	c := f.client(t)
	ok := Query{Start: day(10, 1), End: day(10, 2)}
	bad := map[string]Query{
		"no dates":        {},
		"no end":          {Start: ok.Start},
		"end before":      {Start: ok.End, End: ok.Start},
		"negative limit":  {Start: ok.Start, End: ok.End, RowLimit: -1},
		"limit too big":   {Start: ok.Start, End: ok.End, RowLimit: MaxRowLimit + 1},
		"negative offset": {Start: ok.Start, End: ok.End, StartRow: -1},
		"empty filter":    {Start: ok.Start, End: ok.End, Filters: []Filter{{Dimension: DimensionPage}}},
	}
	for name, q := range bad {
		if _, err := c.Query(context.Background(), "sc-domain:example.com", q); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
	if _, err := c.Query(context.Background(), "", ok); err == nil {
		t.Error("empty site URL: want an error")
	}
	if len(f.paths) != 0 {
		t.Errorf("invalid queries reached the API: %v", f.paths)
	}
}

func TestQueryUsesOnlyTheDate(t *testing.T) {
	var sent map[string]any
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		queryPath: func(w http.ResponseWriter, body []byte) {
			if err := json.Unmarshal(body, &sent); err != nil {
				t.Error(err)
			}
			if _, err := w.Write([]byte(`{}`)); err != nil {
				t.Error(err)
			}
		},
	})
	late := time.Date(2026, 10, 1, 23, 59, 0, 0, time.FixedZone("x", -8*3600))
	if _, err := f.client(t).Query(context.Background(), "sc-domain:example.com", Query{Start: late, End: late}); err != nil {
		t.Fatal(err)
	}
	if sent["startDate"] != "2026-10-01" || sent["endDate"] != "2026-10-01" {
		t.Errorf("dates = %v..%v", sent["startDate"], sent["endDate"])
	}
}

func TestTotals(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		queryPath: reply(`{"rows":[{"clicks":10,"impressions":200,"ctr":0.05,"position":7.2}]}`),
	})
	got, err := f.client(t).Totals(context.Background(), "sc-domain:example.com", day(10, 1), day(10, 7))
	if err != nil {
		t.Fatal(err)
	}
	if got.Clicks != 10 || got.Impressions != 200 || got.Position != 7.2 {
		t.Errorf("totals = %+v", got)
	}
}

func TestTotalsWithoutDataIsZero(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){queryPath: reply(`{}`)})
	got, err := f.client(t).Totals(context.Background(), "sc-domain:example.com", day(10, 1), day(10, 7))
	if err != nil || got.Clicks != 0 || got.Impressions != 0 || got.Keys != nil {
		t.Fatalf("got %+v, %v", got, err)
	}
}

func TestDailyIsOrderedAndParsed(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		queryPath: reply(`{"rows":[
			{"keys":["2026-10-03"],"clicks":3},
			{"keys":["2026-10-01"],"clicks":1},
			{"keys":["2026-10-02"],"clicks":2}]}`),
	})
	got, err := f.client(t).Daily(context.Background(), "sc-domain:example.com", day(10, 1), day(10, 3))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].Clicks != 1 || got[1].Clicks != 2 || got[2].Clicks != 3 {
		t.Fatalf("daily = %+v", got)
	}
	if got[0].Date.Format("2006-01-02") != "2026-10-01" {
		t.Errorf("date = %v", got[0].Date)
	}
}

func TestDailyRejectsAMalformedRow(t *testing.T) {
	for name, body := range map[string]string{
		"no keys":  `{"rows":[{"clicks":1}]}`,
		"bad date": `{"rows":[{"keys":["yesterday"],"clicks":1}]}`,
	} {
		f := newFake(t, map[string]func(http.ResponseWriter, []byte){queryPath: reply(body)})
		if _, err := f.client(t).Daily(context.Background(), "sc-domain:example.com", day(10, 1), day(10, 3)); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}

func TestTop(t *testing.T) {
	var sent map[string]any
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		queryPath: func(w http.ResponseWriter, body []byte) {
			if err := json.Unmarshal(body, &sent); err != nil {
				t.Error(err)
			}
			if _, err := w.Write([]byte(`{"rows":[{"keys":["glossary"],"clicks":9},{"keys":["terms"],"clicks":4}]}`)); err != nil {
				t.Error(err)
			}
		},
	})
	c := f.client(t)
	rows, err := c.Top(context.Background(), "sc-domain:example.com", DimensionQuery, day(10, 1), day(10, 7), 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Keys[0] != "glossary" || sent["rowLimit"] != float64(5) {
		t.Errorf("rows = %+v, request = %v", rows, sent)
	}
	if _, err := c.Top(context.Background(), "sc-domain:example.com", DimensionDate, day(10, 1), day(10, 7), 5); err == nil {
		t.Error("top by date: want an error")
	}
	if _, err := c.Top(context.Background(), "sc-domain:example.com", DimensionPage, day(10, 1), day(10, 7), 0); err == nil {
		t.Error("zero limit: want an error")
	}
}

func TestSiteURLIsEscapedInThePath(t *testing.T) {
	var raw string
	f := newFake(t, nil)
	f.srv.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw = r.URL.EscapedPath()
		if _, err := w.Write([]byte(`{}`)); err != nil {
			t.Error(err)
		}
	})
	if _, err := f.client(t).Totals(context.Background(), "https://www.example.com/", day(10, 1), day(10, 2)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(raw, "sites/https%3A%2F%2Fwww.example.com%2F/searchAnalytics") {
		t.Errorf("escaped path = %q", raw)
	}
}
