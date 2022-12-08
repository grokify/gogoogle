// staticmap assists with https://developers.google.com/maps/documentation/maps-static/start
package staticmap

import (
	"fmt"
	"strings"

	"github.com/grokify/mogo/data/location"
	"github.com/grokify/mogo/errors/errorsutil"
	"google.golang.org/genproto/googleapis/type/latlng"
)

const (
	SizeTiny  = "tiny"
	SizeMid   = "mid"
	SizeSmall = "small"

	ColorBlack  = "black"
	ColorBlue   = "blue"
	ColorBrown  = "brown"
	ColorGray   = "gray"
	ColorGreen  = "green"
	ColorOrange = "orange"
	ColorRed    = "red"
	ColorPurple = "purple"
	ColorWhite  = "white"
	ColorYellow = "yellow"
	// black, brown, green, purple, yellow, blue, gray, orange, red, white
)

func MarkersLatLngs(latlngs ...latlng.LatLng) Markers {
	m := Markers{
		LatLngs: latlngs}
	return m
}

type Markers struct {
	Size            string
	Color           string
	Label           string
	LatLngPrecision uint
	LatLngs         []latlng.LatLng
}

func (m *Markers) String(latLngPrecision uint) string {
	m.trimSpace()
	parts := []string{}
	if len(m.Size) > 0 {
		parts = append(parts, "size:"+m.Size)
	}
	if len(m.Color) > 0 {
		parts = append(parts, "color:"+m.Color)
	}
	for _, ll := range m.LatLngs {
		parts = append(parts, location.LatLngString(&ll, ",", latLngPrecision))
	}
	return strings.Join(parts, "|")
}

func (m *Markers) trimSpace() {
	m.Color = strings.TrimSpace(m.Color)
	m.Label = strings.TrimSpace(m.Label)
	m.Size = strings.TrimSpace(m.Size)
}

type MarkersSlice []Markers

type MarkersMatrix [][]Markers

func (mm MarkersMatrix) WriteFilesPNG(filenames []string, sm StaticMap, key string) error {
	if len(filenames) != len(mm) {
		return fmt.Errorf("filename and markers count mismatch")
	}
	for i, filename := range filenames {
		sm.MarkersList = mm[i]
		err := sm.WriteFilePNG(filename, key)
		if err != nil {
			return errorsutil.Wrapf(err, "filename (%s)", filename)
		}
	}
	return nil
}
