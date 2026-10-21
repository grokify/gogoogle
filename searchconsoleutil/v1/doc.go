// Package searchconsoleutil reads Google Search Console: the properties a
// credential can see, search performance, sitemaps and URL inspection.
//
// It wraps google.golang.org/api/searchconsole/v1 with typed results,
// date handling, input validation and errors that can be tested with
// errors.Is. Everything it calls is read-only and uses the
// webmasters.readonly scope.
//
// Create a Client from a service account key with
// [NewFromServiceAccountKey], or from any authorized *http.Client with [New].
// The service account's email must be added as a user on each Search
// Console property it should read.
//
// Search Console reports dates in Pacific time and holds data back for a few
// days. Where a method takes a time.Time, only its calendar date is used.
package searchconsoleutil
