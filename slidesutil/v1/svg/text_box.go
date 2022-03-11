package slidesutil

// EXPERIMENTAL - NON-WORKING

import (
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
	// ahaslides "github.com/grokify/go-aha/ahaslides"
	"google.golang.org/api/slides/v1"
)

/*
func SVGStart(canvas *svg.SVG, canvasInfo ahaslides.CanvasFloat64) {
	canvas.Start(
		int(canvasInfo.MaxX-canvasInfo.MinX),
		int(canvasInfo.MaxY-canvasInfo.MinY))
}
*/

func SVGAddTextBox(canvas *svg.SVG, tbox CreateShapeTextBoxRequestInfo) {
	// Example: https://www.w3.org/wiki/SVG_Links#Embedding_external_resources_in_an_SVG_document
	tbox.URL = strings.TrimSpace(tbox.URL)
	tbox.Text = strings.TrimSpace(tbox.Text)
	if len(tbox.URL) > 0 {
		canvas.Link(tbox.URL, tbox.Text)
	}
	rectStyles := []string{}
	if len(tbox.BackgroundColorHex) > 0 {
		rectStyles = append(rectStyles, "fill:"+tbox.BackgroundColorHex)
	}
	rectStylesStr := strings.Join(rectStyles, ";")
	canvas.Rect(int(tbox.LocationX), int(tbox.LocationY),
		int(tbox.Width), int(tbox.Height), rectStylesStr)
	padding := 2
	if len(tbox.Text) > 0 {
		stylesText := []string{}
		tbox.ForegroundColorHex = strings.TrimSpace(tbox.ForegroundColorHex)
		if len(tbox.ForegroundColorHex) > 0 {
			stylesText = append(stylesText, "fill:"+tbox.ForegroundColorHex)
		}
		if tbox.FontSize > 0 {
			style := "font-size:" + strconv.Itoa((int(tbox.FontSize))) + strings.TrimSpace(tbox.FontSizeUnit)
			stylesText = append(stylesText, style)
		}
		styleTextStr := strings.Join(stylesText, ";")
		canvas.Text(
			int(tbox.LocationX)+padding,
			int(tbox.LocationY)+padding,
			tbox.Text,
			styleTextStr, // "text-anchor:middle;font-size:30px;fill:white"
		)
	}
	if len(tbox.URL) > 0 {
		canvas.LinkEnd()
	}
}

type CreateShapeTextBoxRequestInfo struct {
	PageID             string
	ObjectID           string
	Width              float64
	Height             float64
	DimensionUnit      string
	LocationX          float64
	LocationY          float64
	LocationUnit       string
	URL                string
	Text               string
	TextURL            string // not implemented yet
	FontBold           bool
	FontItalic         bool
	FontSize           float64
	FontSizeUnit       string
	ParagraphAlignment string // ALIGNMENT_UNSPECIFIED, START, CENTER, END, JUSTIFIED
	ForegroundColorRgb *slides.RgbColor
	BackgroundColorRgb *slides.RgbColor
	ForegroundColorHex string
	BackgroundColorHex string
}
