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

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/google"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"google.golang.org/api/slides/v1"

	su "github.com/grokify/googleutil/slidesutil/v1"
	suex "github.com/grokify/googleutil/slidesutil/v1/examples"
)

func NewClient(forceNewToken bool) (*http.Client, error) {
	conf, err := google.ConfigFromEnv(google.ClientSecretEnv,
		[]string{slides.DriveScope, slides.PresentationsScope})
	if err != nil {
		return nil, err
	}

	tokenFile := "slides.googleapis.com-go-quickstart.json"
	tokenStore, err := goauth.NewTokenStoreFileDefault(tokenFile, true, 0700)
	if err != nil {
		return nil, err
	}

	return goauth.NewClientWebTokenStore(context.Background(), conf, tokenStore, forceNewToken, "mystate")
}

type CreateLineRequestInfo struct {
	LineID        string
	PageID        string
	LineCategory  string // STRAIGHT
	Height        float64
	Width         float64
	DimensionUnit string
	LocationX     float64
	LocationY     float64
	DashStyle     string
	Weight        float64
	ColorHex      string
}

func (info *CreateLineRequestInfo) Requests() ([]*slides.Request, error) {
	reqs := []*slides.Request{
		{
			CreateLine: &slides.CreateLineRequest{
				ObjectId:     info.LineID,
				LineCategory: info.LineCategory,
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: info.PageID,
					Size: &slides.Size{
						Height: &slides.Dimension{Magnitude: info.Height, Unit: info.DimensionUnit},
						Width:  &slides.Dimension{Magnitude: info.Width, Unit: info.DimensionUnit},
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: info.LocationX,
						TranslateY: info.LocationY,
						Unit:       info.DimensionUnit,
					},
				},
			},
		},
	}
	if len(info.ColorHex) > 0 || len(info.DashStyle) > 0 || info.Weight > 0 {
		req := &slides.Request{
			UpdateLineProperties: &slides.UpdateLinePropertiesRequest{
				ObjectId:       info.LineID,
				Fields:         "*",
				LineProperties: &slides.LineProperties{},
			},
		}
		if len(info.ColorHex) > 0 {
			c, err := su.ParseRgbColorHex(info.ColorHex)
			if err != nil {
				return reqs, err
			}
			req.UpdateLineProperties.LineProperties.LineFill = &slides.LineFill{
				SolidFill: &slides.SolidFill{
					Color: &slides.OpaqueColor{
						RgbColor: c,
					},
				},
			}
		}
		if len(info.DashStyle) > 0 {
			req.UpdateLineProperties.LineProperties.DashStyle = info.DashStyle
		}
		if info.Weight > 0 {
			req.UpdateLineProperties.LineProperties.Weight = &slides.Dimension{
				Magnitude: info.Weight,
				Unit:      "PT",
			}
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func CreateLineRequest(pageId string) []*slides.Request {
	lineId := "myLineId"
	reqs := []*slides.Request{
		{
			CreateLine: &slides.CreateLineRequest{
				ObjectId:     lineId,
				LineCategory: "STRAIGHT",
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: pageId,
					Size: &slides.Size{
						Height: &slides.Dimension{Magnitude: 500.0, Unit: "PT"},
						Width:  &slides.Dimension{Magnitude: 1.0, Unit: "PT"},
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: 350.0,
						TranslateY: 100.0,
						Unit:       "PT",
					},
				},
			},
		},
		{
			UpdateLineProperties: &slides.UpdateLinePropertiesRequest{
				ObjectId: lineId,
				Fields:   "*",
				LineProperties: &slides.LineProperties{
					DashStyle: "DASH",
					LineFill: &slides.LineFill{
						SolidFill: &slides.SolidFill{
							Color: &slides.OpaqueColor{
								RgbColor: su.MustParseRgbColorHex("#ff8800"),
							},
						},
					},
					Weight: &slides.Dimension{
						Magnitude: 1.0,
						Unit:      "PT",
					},
				},
			},
		},
	}
	req := &slides.Request{
		UpdatePageElementTransform: &slides.UpdatePageElementTransformRequest{
			ApplyMode: "ABSOLUTE",
			ObjectId:  lineId,
			Transform: &slides.AffineTransform{
				TranslateX: 10,
				TranslateY: 10,
				Unit:       "PT",
			},
		},
	}
	fmtutil.PrintJSON(req)
	return reqs
}

/*
	UpdateLineProperties *UpdateLinePropertiesRequest
*/

// https://productforums.google.com/forum/#!topic/docs/QWeGY_k9hJw
func BorderRequests(objectId string) []*slides.Request {
	req := &slides.Request{
		UpdateTableBorderProperties: &slides.UpdateTableBorderPropertiesRequest{
			ObjectId: objectId,
			Fields:   "*", // "tableBorderFill.solidFill.color"
			TableRange: &slides.TableRange{
				ColumnSpan: 5,
				RowSpan:    8,
			},
			BorderPosition: "ALL",
			TableBorderProperties: &slides.TableBorderProperties{
				TableBorderFill: &slides.TableBorderFill{
					SolidFill: &slides.SolidFill{
						Alpha: 0.0,
						/*Color: &slides.OpaqueColor{
							RgbColor: su.MustParseRgbColorHex("#6f6f6f"),
						},*/
					},
				},
				Weight:    &slides.Dimension{Magnitude: float64(1.0), Unit: "PT"},
				DashStyle: "DASH",
			},
		},
	}

	reqs := []*slides.Request{req}

	return reqs
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

	srv, err := slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides Client %v", err)
	}

	psv := slides.NewPresentationsService(srv)

	pres := &slides.Presentation{Title: "GOLANG TEST PRES #1 TABLE"}
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

	pageID := res.Slides[0].ObjectId
	elementID := "MyTextBox_01"

	requests := []*slides.Request{
		{
			CreateShape: &slides.CreateShapeRequest{
				ObjectId:  elementID,
				ShapeType: "TEXT_BOX",
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: pageID,
					Size: &slides.Size{
						Height: &slides.Dimension{Magnitude: 350, Unit: "PT"},
						Width:  &slides.Dimension{Magnitude: 350, Unit: "PT"},
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: 350.0,
						TranslateY: 100.0,
						Unit:       "PT",
					},
				},
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       elementID,
				InsertionIndex: 0,
				Text:           "New Box Text Inserted!",
			},
		},
	}

	tableID := "MyTable_01"
	rowCount := int64(8)
	columnCount := int64(5)

	requests = []*slides.Request{
		// Create a table
		{
			CreateTable: &slides.CreateTableRequest{
				ObjectId: tableID,
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: pageID,
				},
				Rows:    rowCount,
				Columns: columnCount,
			},
		},
	}

	headings := []string{"foo", "bar", "baz", "qux", "quuz"}
	for i, heading := range headings {
		req := &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId: tableID,
				CellLocation: &slides.TableCellLocation{
					ColumnIndex: int64(i),
					RowIndex:    0,
				},
				Text:           heading,
				InsertionIndex: 0,
			},
		}
		requests = append(requests, req)
	}

	if 1 == 0 {
		borderReqs := BorderRequests(tableID)
		requests = append(requests, borderReqs...)
	}

	if 1 == 0 {
		newReqs, err := su.AlternateRowBgColor(tableID, rowCount, columnCount, "", "#ededed")
		if err != nil {
			panic(err)
		}
		requests = append(requests, newReqs...)
	}

	if 1 == 0 {
		headerStyle := su.UpdateTextStyle{
			ObjectID:           tableID,
			RowIndex:           0,
			ColumnIndex:        0,
			ForegroundColorHex: "#ff8800",
			FontSizeMagnitude:  18,
			FontSizeUnit:       "PT",
		}
		hreqs, err := headerStyle.RequestsColumnSpan(columnCount)
		if err != nil {
			panic(err)
		}
		requests = append(requests, hreqs...)
	}

	if 1 == 0 {
		//lineReqs := CreateLineRequest(pageID)
		lineReqs := suex.LineExampleRequests(pageID)
		requests = append(requests, lineReqs...)
	}

	/*
	   type CreateLineRequestInfo struct {
	   	LineIds       string
	   	PageId        string
	   	LineCategory  string // STRAIGHT
	   	Height        float64
	   	Width         float64
	   	DimensionUnit string
	   	LocationX     float64
	   	LocationY     float64
	   	DashStyle     string
	   	Weight        float64
	   	ColorHex      string
	   }
	*/
	if 1 == 1 {
		for i := 0; i < 5; i++ {
			locX := float64(i) * 100
			lineID := fmt.Sprintf("MYVertLine%03d", i)
			lineInfo := CreateLineRequestInfo{
				PageID:        pageID,
				LineID:        lineID,
				ColorHex:      "#6f6f6f",
				LineCategory:  "STRAIGHT",
				Height:        500.0,
				Width:         1.0,
				DimensionUnit: "PT",
				LocationX:     locX,
				LocationY:     100.0,
				DashStyle:     "DASH",
				Weight:        1.0,
			}
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
