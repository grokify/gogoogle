package slidesutil

import (
	"github.com/pkg/errors"
	"google.golang.org/api/slides/v1"
)

// CreateSlideImageRequestsSidebarRight creates a slide using a main image
// as the body given a PresentationID, title, and imageURL.
func CreateSlideImageRequestsSidebarRight(slideID, imageID, imageURL, sidebarText string) ([]*slides.Request, error) {
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
			return requests, errors.Wrap(err, "")
		}
		requests = append(requests, tbreq...)
	}
	return requests, nil
}
