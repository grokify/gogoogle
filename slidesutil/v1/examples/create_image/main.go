// Go example that covers:
// Quickstart: https://developers.google.com/slides/quickstart/go
// Basic writing: adding a text box to slide: https://developers.google.com/slides/samples/writing
// Using SDK: https://github.com/google/google-api-go-client/blob/master/slides/v1/slides-gen.go
// Creating and Managing Presentations https://developers.google.com/slides/how-tos/presentations
// Adding Shapes and Text to a Slide: https://developers.google.com/slides/how-tos/add-shape#example
// Adding Image to a Slide: https://developers.google.com/slides/how-tos/add-image
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grokify/gogoogle/auth"
	"github.com/grokify/gogoogle/slidesutil/v1"
	"google.golang.org/api/slides/v1"
)

func main() {
	imageURL := "http://11111111.ngrok.io/logo_google_slides.png"
	imageURL = "http://11111111.ngrok.io/chart.png"
	name := "Test Image " + time.Now().Format(time.RFC3339)
	imageID := "MyImageId_01"

	// Requires Drive and Presentations Scopes
	googHttpClient, err := auth.Setup(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	slidesClient, err := slidesutil.NewSlidesClient(googHttpClient)
	if err != nil {
		log.Fatal(err)
	}

	presentationID, err := slidesClient.CreateEmptyPresentation(name)
	if err != nil {
		err = auth.WrapError(err)
		log.Fatal(err)
	}

	slideID, err := slidesClient.CreateSlideTitleAndBody(presentationID, name)
	if err != nil {
		log.Fatal(err)
	}

	requests := CreateSlideRequests(slideID, imageID, imageURL)
	breq := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}

	resu, err := slidesClient.BatchUpdate(presentationID, breq).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resu.PresentationId)
	fmt.Println("DONE")
}

func CreateSlideRequests(slideID, imageID, imageURL string) []*slides.Request {
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
						Width:  &emu4M},
					Transform: &slides.AffineTransform{
						ScaleX:     1.15,
						ScaleY:     1.15,
						TranslateX: 400000.0,
						TranslateY: -300000.0,
						Unit:       "EMU"},
				},
			},
		},
	}
	return requests
}
