package main

import (
	"fmt"

	"github.com/grokify/gotilla/fmt/fmtutil"

	dlpu "github.com/grokify/googleutil/dlp/v2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

/*

This is based on the following Quickstart

https://github.com/GoogleCloudPlatform/golang-samples/blob/master/dlp/dlp_quickstart/quickstart.go

*/

func main() {
	projectID := "PROJECT_ID"
	input := "Hello World"

	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: []*dlppb.InfoType{
			{
				Name: dlpu.InfoTypeAllBasic,
			},
		},

		MinLikelihood: dlppb.Likelihood_POSSIBLE,

		Limits: &dlppb.InspectConfig_FindingLimits{
			MaxFindingsPerRequest: int32(0),
		},
		IncludeQuote: true}

	req := &dlppb.InspectContentRequest{
		Parent:        "projects/" + projectID,
		InspectConfig: inspectConfig,
		Item:          dlpu.NewContentDataItemSimple(input)}

	fmtutil.PrintJSON(req)

	fmt.Println("DONE")
}
