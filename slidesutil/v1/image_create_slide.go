package slidesutil

import (
	"strings"
	"time"

	"github.com/grokify/mogo/crypto/md5util"
	"github.com/grokify/mogo/errors/errorsutil"
	slides "google.golang.org/api/slides/v1"
)

// CreateSlideImageRequestsSidebarRight creates API batch requests to
// load a main page with optional right sidebar text. `imageID` is
// optional and will be auto-generated if not provided.
func CreateSlideImageRequestsSidebarRight(slideID, imageID, imageURL, sidebarText string) ([]*slides.Request, error) {
	slideID = strings.TrimSpace(slideID)
	imageID = strings.TrimSpace(imageID)
	imageURL = strings.TrimSpace(imageURL)

	requests := []*slides.Request{}
	if len(imageURL) > 0 {
		if len(imageID) == 0 {
			imageID = md5util.MD5Base62(slideID + "." + imageURL + "." + time.Now().Format(time.RFC3339))
		}
		emu4M := slides.Dimension{Magnitude: 6000000, Unit: "EMU"}
		requests = append(requests, &slides.Request{
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
		})
	}
	if len(sidebarText) > 0 {
		textboxes := CreateShapeTextBoxRequestInfo{
			PageId:        slideID,
			ObjectId:      imageID + ":txt",
			Text:          sidebarText,
			Width:         90,
			Height:        50,
			FontSize:      10,
			FontSizeUnit:  UnitPT,
			DimensionUnit: UnitPT,
			LocationX:     580,
			LocationY:     158,
			LocationUnit:  UnitPT,
		}
		tbreq, err := textboxes.Requests()
		if err != nil {
			return requests, errorsutil.Wrap(err, "textboxes.Requests")
		}
		requests = append(requests, tbreq...)
	}
	return requests, nil
}
