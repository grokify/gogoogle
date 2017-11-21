package slidesutil

import (
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
	"google.golang.org/api/slides/v1"
)

var (
	GoogleSlideUnitPoint = "PT"
)

func RgbColorParseHex(hexColor string) (*slides.RgbColor, error) {
	c, err := colorful.Hex(hexColor)
	if err != nil {
		return nil, err
	}
	return &slides.RgbColor{Red: c.R, Green: c.G, Blue: c.B}, nil
}

func RgbColorMustParseHex(hexColor string) *slides.RgbColor {
	c, err := colorful.Hex(hexColor)
	if err != nil {
		panic(`colorful: Hex(` + quote(hexColor) + `): ` + err.Error())
	}
	return &slides.RgbColor{Red: c.R, Green: c.G, Blue: c.B}
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func TextBoxRequestsSimple(pageId, elementId, text string, fgColor, bgColor *slides.RgbColor, width, height, locX, locY float64) []*slides.Request {
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
