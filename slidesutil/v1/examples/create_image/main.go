// Go example that covers:
// Quickstart: https://developers.google.com/slides/quickstart/go
// Basic writing: adding a text box to slide: https://developers.google.com/slides/samples/writing
// Using SDK: https://github.com/google/google-api-go-client/blob/master/slides/v1/slides-gen.go
// Creating and Managing Presentations https://developers.google.com/slides/how-tos/presentations
// Adding Shapes and Text to a Slide: https://developers.google.com/slides/how-tos/add-shape#example
// Adding Image to a Slide: https://developers.google.com/slides/how-tos/add-image
package main

import (
	"fmt"
	"log"

	"github.com/grokify/googleutil/slidesutil/v1"
	slidesutilexamples "github.com/grokify/googleutil/slidesutil/v1/examples"
	"google.golang.org/api/slides/v1"
)

func main() {
	imageURL := "http://11111111.ngrok.io/logo_google_slides.png"
	imageURL = "http://11111111.ngrok.io/chart.png"

	gss, err := slidesutilexamples.Setup()
	if err != nil {
		log.Fatal(err)
	}

	presentationID, err := slidesutil.CreateEmptyPresentation(gss.PresentationsService, "Test Image")
	if err != nil {
		log.Fatal(err)
	}
	slideID, err := slidesutil.CreateSlideTitleOnly(
		gss.SlidesService, gss.PresentationsService,
		presentationID, "Test Image")

	imageID := "MyImageId_01"
	emu4M := slides.Dimension{Magnitude: 6000000, Unit: "EMU"}

	requests := []*slides.Request{
		{
			CreateImage: &slides.CreateImageRequest{
				ObjectId: imageID,
				Url:      imageURL,
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: slideID,
					Size: &slides.Size{
						Height: &emu4M,
						Width:  &emu4M,
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: 100000.0,
						TranslateY: 100000.0,
						Unit:       "EMU",
					},
				},
			},
		},
	}
	breq := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}

	resu, err := gss.PresentationsService.BatchUpdate(presentationID, breq).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resu.PresentationId)
}
