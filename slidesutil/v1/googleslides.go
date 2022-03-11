package slidesutil

import (
	"regexp"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
	"google.golang.org/api/slides/v1"
)

var (
	GoogleSlideUnitPoint = "PT"
	ObjectIDFormat       = `^[a-zA-Z0-9_][a-zA-Z0-9_\-:]*$`
)

/*
panic: googleapi: Error 400: Invalid requests[0].createShape: The object ID (TAGLABELBG-rmp-Glip Growth) should start with a word character [a-zA-Z0-9_] and then followed by any number of the following characters [a-zA-Z0-9_-:]., badRequest
*/

func FormatObjectIDSimple(s string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9_\-:]`).ReplaceAllString(s, "_")
}

func ParseRgbColorHex(hexColor string) (*slides.RgbColor, error) {
	c, err := colorful.Hex(hexColor)
	if err != nil {
		return nil, err
	}
	return &slides.RgbColor{Red: c.R, Green: c.G, Blue: c.B}, nil
}

func MustParseRgbColorHex(hexColor string) *slides.RgbColor {
	c, err := ParseRgbColorHex(hexColor)
	if err != nil {
		panic(`ParseColor: Hex(` + quote(hexColor) + `): ` + err.Error())
	}
	return c
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

type TextBoxInfoSimple struct {
	PageID     string
	ElementID  string
	Text       string
	FgColorHex string
	BgColorHex string
	Width      float64
	Height     float64
	LocationX  float64
	LocationY  float64
}

func AddFgColor(elementID string, fgColor *slides.RgbColor) []*slides.Request {
	reqs := []*slides.Request{
		{
			UpdateTextStyle: &slides.UpdateTextStyleRequest{
				ObjectId: elementID,
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
	}
	return reqs
}
