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

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/google"
	"github.com/grokify/mogo/math/mathutil"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"

	su "github.com/grokify/gogoogle/slidesutil/v1"
	suex "github.com/grokify/gogoogle/slidesutil/v1/examples"
)

func NewClient(forceNewToken bool) (*http.Client, error) {
	conf, err := google.ConfigFromEnv(google.ClientSecretEnv,
		[]string{slides.DriveScope, slides.PresentationsScope})
	if err != nil {
		return nil, err
	}

	tokenFile := "slides.googleapis.com-go-quickstart.json"
	tokenStore, err := authutil.NewTokenStoreFileDefault(tokenFile, true, 0700)
	if err != nil {
		return nil, err
	}

	return authutil.NewClientWebTokenStore(context.Background(), conf, tokenStore, forceNewToken, "mystate")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	forceNewToken := false

	client, err := NewClient(forceNewToken)
	if err != nil {
		log.Fatal("Unable to get Client")
	}

	srv, err := slides.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("unable to retrieve Slides Client %v", err)
	}

	psv := slides.NewPresentationsService(srv)

	pres := &slides.Presentation{Title: "GOLANG TEST PRES #1 TABLE"}
	res, err := psv.Create(pres).Do()
	if err != nil {
		panic(err)
	}

	fmt.Printf("created Presentation with Id %v\n", res.PresentationId)

	for i, slide := range res.Slides {
		fmt.Printf("- Slide #%d id %v contains %d elements.\n", (i + 1),
			slide.ObjectId,
			len(slide.PageElements))
	}

	pageID := res.Slides[0].ObjectId

	requests := []*slides.Request{}

	// Single line
	if 1 == 0 {
		lineReqs := suex.LineExampleRequests(pageID)
		requests = append(requests, lineReqs...)
	}

	// Simple for loop
	if 1 == 0 {
		lineInfo := su.CreateLineRequestInfo{
			PageID:        pageID,
			LineID:        "",
			ColorHex:      "#6f6f6f",
			LineCategory:  "STRAIGHT",
			Height:        300.0,
			Width:         1.0,
			DimensionUnit: "PT",
			LocationX:     0.0,
			LocationY:     100.0,
			DashStyle:     "DASH",
			Weight:        1.0,
		}
		for i := 0; i < 5; i++ {
			lineInfo.LineID = fmt.Sprintf("MYVertLine%03d", i)
			lineInfo.LocationX = float64(i) * 100
			lineReqs, err := lineInfo.Requests()
			if err != nil {
				panic(err)
			}
			requests = append(requests, lineReqs...)
		}
	}

	// Using a range struct.
	boolUseRangeStruct := true
	if boolUseRangeStruct {
		rng := mathutil.RangeFloat64{
			Min:   150.0,
			Max:   700.0,
			Cells: int32(5),
		}
		linePrefix := "VertLines"
		lineInfo := su.CreateLineRequestInfo{
			PageID:        pageID,
			LineID:        "",
			ColorHex:      "#6f6f6f",
			LineCategory:  "STRAIGHT",
			Height:        300.0,
			Width:         1.0,
			DimensionUnit: "PT",
			LocationX:     0.0,
			LocationY:     100.0,
			DashStyle:     "DASH",
			Weight:        1.0,
		}
		for i := int32(0); i < rng.Cells; i++ {
			min, max, err := rng.CellMinMax(i)
			if err != nil {
				panic(err)
			}
			fmt.Printf("IDX %v MIN %v MAX %v\n", i, min, max)
			if i == 0 {
				lineInfo.LineID = fmt.Sprintf("%v%03d%v", linePrefix, i, "start")
				lineInfo.LocationX = min
				lineReqs, err := lineInfo.Requests()
				if err != nil {
					panic(err)
				}
				requests = append(requests, lineReqs...)
			}
			lineInfo.LineID = fmt.Sprintf("%v%03d%v", linePrefix, i, "end")
			lineInfo.LocationX = max
			lineReqs, err := lineInfo.Requests()
			if err != nil {
				panic(err)
			}
			requests = append(requests, lineReqs...)
		}
	}

	breq := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}

	resu, err := psv.BatchUpdate(res.PresentationId, breq).Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(resu.PresentationId)
}
