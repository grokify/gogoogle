package slidesutil

import (
	colorful "github.com/lucasb-eyer/go-colorful"
	slides "google.golang.org/api/slides/v1"
)

func OptionalColorParseHex(hexColorStr string) (*slides.OptionalColor, error) {
	rgb, err := RgbColorParseHex(hexColorStr)
	if err != nil {
		return nil, err
	}
	c := &slides.OptionalColor{
		OpaqueColor: &slides.OpaqueColor{
			RgbColor: rgb,
		},
	}
	return c, nil
}

func RgbColorParseHex(hexColorStr string) (*slides.RgbColor, error) {
	col, err := colorful.Hex(hexColorStr)
	if err != nil {
		return nil, err
	}
	return &slides.RgbColor{
		Blue:  col.B,
		Green: col.G,
		Red:   col.R,
	}, nil
}
