# Google Search Console

The `searchconsoleutil/v1` package reads Google Search Console: the properties
a credential can see, search performance, sitemaps and URL inspection. It is
read-only and requests only the `webmasters.readonly` scope.

## Connect with a service account

A service account lets a server read Search Console without a person signing
in.

1. In Google Cloud, enable the **Google Search Console API** and create a
   service account and a JSON key for it.
2. In Search Console, open **Settings**, then **Users and permissions**, and
   add the service account's email as a user on each property it should read.
   The lowest permission that can read is enough.
3. Create the client from the key:

```go
import searchconsole "github.com/grokify/gogoogle/searchconsoleutil/v1"

key, err := os.ReadFile("key.json")
if err != nil {
    return err
}
email, _ := searchconsole.ServiceAccountEmail(key) // the address to add in step 2
c, err := searchconsole.NewFromServiceAccountKey(ctx, key)
```

The context is used each time the access token is refreshed, so pass one that
lasts as long as the client. `NewFromServiceAccountKey` accepts only service
account keys. To use credentials from another source, authorize an
`*http.Client` for `searchconsole.ScopeReadonly` and pass it to
`searchconsole.New`.

## Properties

```go
sites, err := c.Sites(ctx)
for _, s := range sites {
    if s.CanRead() && s.MatchesDomain("example.com") {
        fmt.Println(s.URL, s.PermissionLevel, s.IsDomainProperty())
    }
}
```

A property is either a domain property, `sc-domain:example.com`, covering
every subdomain and protocol, or a URL-prefix property such as
`https://www.example.com/`. The `URL` is the identifier the other methods
take. `MatchesDomain` is true when the property's host is the domain, a
subdomain of it, or a parent domain that covers it.

## Search performance

```go
end := time.Now().AddDate(0, 0, -3) // Search Console holds data back for days
start := end.AddDate(0, 0, -27)

totals, err := c.Totals(ctx, site, start, end)       // clicks, impressions, CTR, position
daily, err := c.Daily(ctx, site, start, end)         // one row per day, oldest first
queries, err := c.Top(ctx, site, searchconsole.DimensionQuery, start, end, 25)
pages, err := c.Top(ctx, site, searchconsole.DimensionPage, start, end, 25)
```

`Query` is the general form, with dimensions, filters, a search type, a data
state and paging:

```go
res, err := c.Query(ctx, site, searchconsole.Query{
    Start: start, End: end,
    Dimensions: []searchconsole.Dimension{searchconsole.DimensionQuery, searchconsole.DimensionDevice},
    Filters: []searchconsole.Filter{{
        Dimension: searchconsole.DimensionPage, Operator: searchconsole.FilterContains, Expression: "/blog/",
    }},
    RowLimit: 100,
})
fmt.Println(res.FirstIncompleteDate) // data from this date may still change
```

Only the calendar date of a `time.Time` is used. Search Console reports dates
in Pacific time.

## Sitemaps and URL inspection

```go
maps, err := c.Sitemaps(ctx, site)

in, err := c.InspectURL(ctx, site, "https://example.com/pricing")
fmt.Println(in.Indexed(), in.CoverageState, in.LastCrawl, in.GoogleCanonical)
```

URL inspection is limited per property per day, so cache what it returns.

## Errors

Failed API calls wrap both a sentinel and the underlying `*googleapi.Error`:

| Error | Meaning |
|---|---|
| `ErrUnauthorized` | Google rejected the credentials |
| `ErrForbidden` | The credentials cannot access the property, for example the service account was not added to it |
| `ErrNotFound` | The property, sitemap or URL does not exist for the credentials |
| `ErrQuotaExceeded` | A rate or daily quota was reached |

```go
if errors.Is(err, searchconsole.ErrForbidden) {
    // ask the owner to add the service account to the property
}
```
