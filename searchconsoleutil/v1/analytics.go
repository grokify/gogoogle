package searchconsoleutil

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	searchconsole "google.golang.org/api/searchconsole/v1"
)

// Dimension is a way to group search analytics rows.
type Dimension string

// Dimensions that rows can be grouped by.
const (
	DimensionDate             Dimension = "date"
	DimensionQuery            Dimension = "query"
	DimensionPage             Dimension = "page"
	DimensionCountry          Dimension = "country"
	DimensionDevice           Dimension = "device"
	DimensionSearchAppearance Dimension = "searchAppearance"
)

// SearchType is the kind of search results to report on.
type SearchType string

// Search types. The empty value means web.
const (
	SearchTypeWeb      SearchType = "web"
	SearchTypeImage    SearchType = "image"
	SearchTypeVideo    SearchType = "video"
	SearchTypeNews     SearchType = "news"
	SearchTypeDiscover SearchType = "discover"
)

// DataState says whether fresh, incomplete data is included.
type DataState string

// Data states. The empty value means final.
const (
	// DataStateFinal reports only data that is complete.
	DataStateFinal DataState = "final"
	// DataStateAll includes fresh data that may still change.
	DataStateAll DataState = "all"
)

// FilterOperator compares a dimension to an expression.
type FilterOperator string

// Filter operators.
const (
	FilterEquals         FilterOperator = "equals"
	FilterNotEquals      FilterOperator = "notEquals"
	FilterContains       FilterOperator = "contains"
	FilterNotContains    FilterOperator = "notContains"
	FilterIncludingRegex FilterOperator = "includingRegex"
	FilterExcludingRegex FilterOperator = "excludingRegex"
)

// MaxRowLimit is the most rows one request can return.
const MaxRowLimit = 25000

// dateLayout is how Search Console writes a date.
const dateLayout = "2006-01-02"

// Filter restricts a query to rows whose dimension matches.
type Filter struct {
	Dimension  Dimension
	Operator   FilterOperator
	Expression string
}

// Query describes a search analytics request.
type Query struct {
	// Start and End bound the dates, inclusive. Only their calendar dates
	// are used. Both are required.
	Start, End time.Time
	// Dimensions group the rows. With none, the result is one row of totals.
	Dimensions []Dimension
	SearchType SearchType
	DataState  DataState
	// Filters must all match.
	Filters []Filter
	// RowLimit is the most rows to return, up to [MaxRowLimit]. Zero uses
	// the API's default of 1000.
	RowLimit int
	// StartRow skips this many rows, for paging through a large result.
	StartRow int
}

// Row is one row of search analytics.
type Row struct {
	// Keys holds the value of each requested dimension, in order.
	Keys        []string
	Clicks      float64
	Impressions float64
	// CTR is the click-through rate from 0 to 1.
	CTR float64
	// Position is the average position in search results, where 1 is the top.
	Position float64
}

// Result is the answer to a Query.
type Result struct {
	Rows []Row
	// FirstIncompleteDate is the first date whose data may still change.
	// It is zero when the response gave none.
	FirstIncompleteDate time.Time
}

// DailyRow is one day's totals.
type DailyRow struct {
	Date time.Time
	Row
}

// validate checks the query before any request is made.
func (q Query) validate() error {
	switch {
	case q.Start.IsZero() || q.End.IsZero():
		return errors.New("searchconsole: query needs a start and an end date")
	case dayOf(q.End).Before(dayOf(q.Start)):
		return errors.New("searchconsole: query end date is before its start date")
	case q.RowLimit < 0 || q.RowLimit > MaxRowLimit:
		return fmt.Errorf("searchconsole: row limit must be between 0 and %d", MaxRowLimit)
	case q.StartRow < 0:
		return errors.New("searchconsole: start row cannot be negative")
	}
	for _, f := range q.Filters {
		if f.Dimension == "" || f.Expression == "" {
			return errors.New("searchconsole: a filter needs a dimension and an expression")
		}
	}
	return nil
}

// dayOf drops the time of day.
func dayOf(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// Query runs a search analytics query against the property siteURL.
func (c *Client) Query(ctx context.Context, siteURL string, q Query) (*Result, error) {
	if siteURL == "" {
		return nil, errors.New("searchconsole: site URL cannot be empty")
	}
	if err := q.validate(); err != nil {
		return nil, err
	}
	req := &searchconsole.SearchAnalyticsQueryRequest{
		StartDate:  dayOf(q.Start).Format(dateLayout),
		EndDate:    dayOf(q.End).Format(dateLayout),
		SearchType: string(q.SearchType),
		DataState:  string(q.DataState),
		RowLimit:   int64(q.RowLimit),
		StartRow:   int64(q.StartRow),
	}
	for _, d := range q.Dimensions {
		req.Dimensions = append(req.Dimensions, string(d))
	}
	if len(q.Filters) > 0 {
		g := &searchconsole.ApiDimensionFilterGroup{GroupType: "and"}
		for _, f := range q.Filters {
			g.Filters = append(g.Filters, &searchconsole.ApiDimensionFilter{
				Dimension: string(f.Dimension), Operator: string(f.Operator), Expression: f.Expression,
			})
		}
		req.DimensionFilterGroups = []*searchconsole.ApiDimensionFilterGroup{g}
	}
	resp, err := c.svc.Searchanalytics.Query(siteURL, req).Context(ctx).Do()
	if err != nil {
		return nil, wrap(err)
	}
	out := &Result{Rows: make([]Row, 0, len(resp.Rows))}
	for _, r := range resp.Rows {
		out.Rows = append(out.Rows, Row{
			Keys: r.Keys, Clicks: r.Clicks, Impressions: r.Impressions, CTR: r.Ctr, Position: r.Position,
		})
	}
	if resp.Metadata != nil && resp.Metadata.FirstIncompleteDate != "" {
		t, err := time.Parse(dateLayout, resp.Metadata.FirstIncompleteDate)
		if err != nil {
			return nil, fmt.Errorf("searchconsole: read first incomplete date: %w", err)
		}
		out.FirstIncompleteDate = t
	}
	return out, nil
}

// Totals returns the clicks, impressions, click-through rate and position
// for the whole period. A period with no data returns a zero Row.
func (c *Client) Totals(ctx context.Context, siteURL string, start, end time.Time) (Row, error) {
	res, err := c.Query(ctx, siteURL, Query{Start: start, End: end})
	if err != nil {
		return Row{}, err
	}
	if len(res.Rows) == 0 {
		return Row{}, nil
	}
	return res.Rows[0], nil
}

// Daily returns one row per day that has data, oldest first.
func (c *Client) Daily(ctx context.Context, siteURL string, start, end time.Time) ([]DailyRow, error) {
	res, err := c.Query(ctx, siteURL, Query{
		Start: start, End: end, Dimensions: []Dimension{DimensionDate}, RowLimit: MaxRowLimit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]DailyRow, 0, len(res.Rows))
	for _, r := range res.Rows {
		if len(r.Keys) != 1 {
			return nil, fmt.Errorf("searchconsole: a daily row has %d keys, want 1", len(r.Keys))
		}
		d, err := time.Parse(dateLayout, r.Keys[0])
		if err != nil {
			return nil, fmt.Errorf("searchconsole: read a daily row's date: %w", err)
		}
		out = append(out, DailyRow{Date: d, Row: r})
	}
	// The API sorts by date already; sorting here keeps the promise when it
	// does not.
	sortDaily(out)
	return out, nil
}

// Top returns the rows with the most clicks for one dimension, such as the
// top queries or the top pages. limit is the most rows, up to [MaxRowLimit].
func (c *Client) Top(ctx context.Context, siteURL string, dim Dimension, start, end time.Time, limit int) ([]Row, error) {
	if dim == "" || dim == DimensionDate {
		return nil, errors.New("searchconsole: top needs a dimension other than date")
	}
	if limit <= 0 {
		return nil, errors.New("searchconsole: top needs a positive limit")
	}
	res, err := c.Query(ctx, siteURL, Query{
		Start: start, End: end, Dimensions: []Dimension{dim}, RowLimit: limit,
	})
	if err != nil {
		return nil, err
	}
	return res.Rows, nil
}

func sortDaily(rows []DailyRow) {
	sort.Slice(rows, func(i, j int) bool { return rows[i].Date.Before(rows[j].Date) })
}
