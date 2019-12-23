package slidesutil

import (
	"google.golang.org/api/slides/v1"
)

const (
	LayoutMainPoint      string = "MAIN_POINT"
	LayoutTitleAndBody   string = "TITLE_AND_BODY"
	RangeTypeUnspecified string = "RANGE_TYPE_UNSPECIFIED"
	RangeTypeFixedRange  string = "FIXED_RANGE"
	RangeTypeStartIndex  string = "FROM_START_INDEX"
	RangeTypeAll         string = "ALL"
)

func CreateSlideRequestLayout(predefinedLayout string) *slides.Request {
	return &slides.Request{
		CreateSlide: &slides.CreateSlideRequest{
			SlideLayoutReference: &slides.LayoutReference{
				PredefinedLayout: predefinedLayout,
			},
		},
	}
}

func InsertTextRequest(objectID, text string) *slides.Request {
	return &slides.Request{
		InsertText: &slides.InsertTextRequest{
			ObjectId: objectID,
			Text:     text,
		},
	}
}

func UpdateTextStyleRequestBold(objectID string, startIdx, endIdx int64) *slides.Request {
	return &slides.Request{
		UpdateTextStyle: &slides.UpdateTextStyleRequest{
			ObjectId:  objectID,
			Style:     &slides.TextStyle{Bold: true},
			TextRange: &slides.Range{Type: RangeTypeFixedRange, StartIndex: startIdx, EndIndex: endIdx},
			Fields:    "bold",
		},
	}
}

func UpdateTextStyleRequestBullet(objectID string, startIdx, endIdx int64) *slides.Request {
	return &slides.Request{
		CreateParagraphBullets: &slides.CreateParagraphBulletsRequest{
			ObjectId:  objectID,
			TextRange: &slides.Range{Type: RangeTypeFixedRange, StartIndex: startIdx, EndIndex: endIdx},
		},
	}
}

func UpdateTextStyleRequestFontSizePT(objectID string, pointSize float64) *slides.Request {
	return UpdateTextStyleRequestFontSize(objectID,
		slides.Dimension{
			Magnitude: pointSize,
			Unit:      "PT"})
}

func UpdateTextStyleRequestFontSize(objectID string, dimension slides.Dimension) *slides.Request {
	return &slides.Request{
		UpdateTextStyle: &slides.UpdateTextStyleRequest{
			ObjectId: objectID,
			Style:    &slides.TextStyle{FontSize: &dimension},
			Fields:   "fontSize",
		},
	}
}

func UpdateParagraphStyleRequestLineSpacing(objectID string, lineSpacing float64) *slides.Request {
	return &slides.Request{
		UpdateParagraphStyle: &slides.UpdateParagraphStyleRequest{
			ObjectId: objectID,
			Style: &slides.ParagraphStyle{
				LineSpacing: lineSpacing,
			},
			Fields: "lineSpacing",
		},
	}
}

/*
func OptionalColorHex() {
	c := &slides.OptionalColor{
		OpaqueColor: &slides.OpaqueColor{
			RgbColor: &slides.RgbColor{},
		},
	}
}
*/
func UpdateTextStyleRequestLinkURL(objectID, url string, textRange *slides.Range, underlineLinks bool) *slides.Request {
	optionalColor, err := OptionalColorParseHex("#666666")
	if err != nil {
		panic(err)
	}
	return &slides.Request{
		UpdateTextStyle: &slides.UpdateTextStyleRequest{
			ObjectId:  objectID,
			TextRange: textRange,
			Style: &slides.TextStyle{
				ForegroundColor: optionalColor,
				Link:            &slides.Link{Url: url},
				Underline:       underlineLinks},
			Fields: "link,underline,foregroundColor",
		},
	}
}
