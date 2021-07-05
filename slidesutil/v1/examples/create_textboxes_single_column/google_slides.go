// Go example that covers:
// Quickstart: https://developers.google.com/slides/quickstart/go
// Basic writing: adding a text box to slide: https://developers.google.com/slides/samples/writing
// Using SDK: https://github.com/google/google-api-go-client/blob/master/slides/v1/slides-gen.go
// Creating and Managing Presentations https://developers.google.com/slides/how-tos/presentations
// Adding Shapes and Text to a Slide: https://developers.google.com/slides/how-tos/add-shape#example
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/grokify/googleutil/slidesutil/v1"
	"github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/google"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"google.golang.org/api/slides/v1"
)

var (
	GoogleSlideUnitPoint = "PT"
)

func NewClient(forceNewToken bool) (*http.Client, error) {
	conf, err := google.ConfigFromEnv(google.ClientSecretEnv,
		[]string{slides.DriveScope, slides.PresentationsScope})
	if err != nil {
		return nil, err
	}

	tokenFile := "slides.googleapis.com-go-quickstart.json"
	tokenStore, err := oauth2more.NewTokenStoreFileDefault(tokenFile, true, 0700)
	if err != nil {
		return nil, err
	}

	return oauth2more.NewClientWebTokenStore(context.Background(), conf, tokenStore, forceNewToken, "mystate")
}

type CreateShapeTextBoxRequestInfo struct {
	PageId             string
	ObjectId           string
	Width              float64
	Height             float64
	DimensionUnits     string
	LocationX          float64
	LocationY          float64
	LocationUnits      string
	Text               string
	ForegroundColorHex string
	BackgroundColorHex string
}

func (info *CreateShapeTextBoxRequestInfo) Requests() ([]*slides.Request, error) {
	requests := []*slides.Request{
		{
			CreateShape: &slides.CreateShapeRequest{
				ObjectId:  info.ObjectId,
				ShapeType: "TEXT_BOX",
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: info.PageId,
					Size: &slides.Size{
						Width:  &slides.Dimension{Magnitude: info.Width, Unit: info.DimensionUnits},
						Height: &slides.Dimension{Magnitude: info.Height, Unit: info.DimensionUnits},
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: info.LocationX,
						TranslateY: info.LocationY,
						Unit:       info.LocationUnits,
					},
				},
			},
		},
	}
	if len(info.Text) > 0 {
		requests = append(requests, &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       info.ObjectId,
				InsertionIndex: 0,
				Text:           info.Text,
			},
		})
	}

	if len(info.ForegroundColorHex) > 0 {

	}
	/*
		{
		UpdateTextStyle: &slides.UpdateTextStyleRequest{
			ObjectId: elementId,
			Fields:   "*",
			Style: &slides.TextStyle{
				FontSize: &slides.Dimension{
					Magnitude: 10.0,
					Unit:      GoogleSlideUnitPoint,
				},
				ForegroundColor: &slides.OptionalColor{
					OpaqueColor: &slides.OpaqueColor{
						RgbColor: fgColor,
					},
				},
			},
		},
	},*/

	return requests, nil
}

func TextBoxRequests(pageId, elementId, text string, fgColor, bgColor *slides.RgbColor, width, height, locX, locY float64) []*slides.Request {
	return []*slides.Request{
		{
			CreateShape: &slides.CreateShapeRequest{
				ObjectId:  elementId,
				ShapeType: "TEXT_BOX",
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: pageId,
					Size: &slides.Size{
						Width:  &slides.Dimension{Magnitude: width, Unit: GoogleSlideUnitPoint},
						Height: &slides.Dimension{Magnitude: height, Unit: GoogleSlideUnitPoint},
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: locX,
						TranslateY: locY,
						Unit:       GoogleSlideUnitPoint,
					},
				},
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       elementId,
				InsertionIndex: 0,
				Text:           text,
			},
		},
		{
			UpdateTextStyle: &slides.UpdateTextStyleRequest{
				ObjectId: elementId,
				Fields:   "*",
				Style: &slides.TextStyle{
					FontSize: &slides.Dimension{
						Magnitude: 10.0,
						Unit:      GoogleSlideUnitPoint,
					},
					ForegroundColor: &slides.OptionalColor{
						OpaqueColor: &slides.OpaqueColor{
							RgbColor: fgColor,
						},
					},
				},
			},
		},
		{
			UpdateShapeProperties: &slides.UpdateShapePropertiesRequest{
				ObjectId: elementId,
				Fields:   "shapeBackgroundFill.solidFill.color",
				ShapeProperties: &slides.ShapeProperties{
					ShapeBackgroundFill: &slides.ShapeBackgroundFill{
						SolidFill: &slides.SolidFill{
							Color: &slides.OpaqueColor{
								RgbColor: bgColor,
							},
						},
					},
				},
			},
		},
	}
}

func main() {
	forceNewToken := false

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	client, err := NewClient(forceNewToken)
	if err != nil {
		log.Fatal("Unable to get Client")
	}

	srv, err := slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides Client %v", err)
	}

	psv := slides.NewPresentationsService(srv)

	t := time.Now().UTC()
	slideName := fmt.Sprintf("GOLANG TEST PRES %v\n", t.Format(time.RFC3339))
	fmt.Printf("Slide Name: %v", slideName)

	pres := &slides.Presentation{Title: slideName}
	res, err := psv.Create(pres).Do()
	if err != nil {
		panic(err)
	}

	fmt.Printf("CREATED Presentation with Id %v\n", res.PresentationId)

	for i, slide := range res.Slides {
		fmt.Printf("- Slide #%d id %v contains %d elements.\n", (i + 1),
			slide.ObjectId,
			len(slide.PageElements))
	}

	pageId := res.Slides[0].ObjectId
	requests := []*slides.Request{}

	//fgColor, err := colorutil.GoogleSlidesRgbColorParseHex("#ffffff")
	fgColor, err := slidesutil.ParseRgbColorHex("#ffffff")
	if err != nil {
		panic(err)
	}
	//bgColor, err := colorutil.GoogleSlidesRgbColorParseHex("#4688f1")
	bgColor, err := slidesutil.ParseRgbColorHex("#4688f1")
	if err != nil {
		panic(err)
	}

	items := []string{"Item #1", "Item #2"}
	locX := 350.0
	locY := 50.0
	boxWidth := 130.0
	boxHeight := 25.0
	locYHeight := boxHeight + 5.0
	for i, itemText := range items {
		elementId := fmt.Sprintf("item%v", i)
		locYThis := locY + locYHeight*float64(i)
		requests = append(requests, TextBoxRequests(
			pageId, elementId, itemText, fgColor, bgColor,
			boxWidth, boxHeight, locX, locYThis)...)
	}

	breq := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}

	resu, err := psv.BatchUpdate(res.PresentationId, breq).Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(resu.PresentationId)

	fmt.Println("DONE")
}
