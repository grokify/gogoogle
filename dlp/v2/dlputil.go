package dlp

import (
	// dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// For all InvoTypes, see: https://cloud.google.com/dlp/docs/infotypes-reference

const (
	InfoTypeAllBasic   = "ALL_BASIC"
	InfoTypePersonName = "PERSON_NAME"
	InfoTypeUSState    = "US_STATE"
)

/* Likelihook:

https://godoc.org/google.golang.org/genproto/googleapis/privacy/dlp/v2#Likelihood

const (
	// Default value; same as POSSIBLE.
	Likelihood_LIKELIHOOD_UNSPECIFIED Likelihood = 0
	// Few matching elements.
	Likelihood_VERY_UNLIKELY Likelihood = 1
	Likelihood_UNLIKELY      Likelihood = 2
	// Some matching elements.
	Likelihood_POSSIBLE Likelihood = 3
	Likelihood_LIKELY   Likelihood = 4
	// Many matching elements.
	Likelihood_VERY_LIKELY Likelihood = 5
)
*/

func NewContentDataItemSimple(input string) *dlppb.ContentItem {
	return &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: input,
		}}
}
