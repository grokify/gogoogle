package bigqueryutil

import (
	"context"

	"cloud.google.com/go/bigquery"
	iu "github.com/grokify/gotilla/type/interfaceutil"
)

const (
	MaxInsertsPerOperation = 10000
)

// UploadItems uses Google's streaming upload: https://cloud.google.com/bigquery/streaming-data-into-bigquery
func UploadItems(ctx context.Context, client *bigquery.Client, datasetID, tableID string, items []interface{}) []error {
	u := client.Dataset(datasetID).Table(tableID).Uploader()
	items2 := iu.SplitSliceInterface(items, MaxInsertsPerOperation)

	errs := []error{}

	for _, items1 := range items2 {
		err := u.Put(ctx, items1)
		if err != nil {
			if multiError, ok := err.(bigquery.PutMultiError); ok {
				for _, ex := range multiError {
					for _, ey := range ex.Errors {
						errs = append(errs, ey)
					}
				}
			} else {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
