package slidesutil

import (
	"fmt"

	"github.com/grokify/gocharts/v2/data/roadmap"
	"google.golang.org/api/slides/v1"
)

type SlideCanvasInfo struct {
	BoxFgColor      *slides.RgbColor
	BoxBgColor      *slides.RgbColor
	BoxHeight       float64
	BoxMarginBottom float64
	Canvas          CanvasFloat64
}

func DefaultSlideCanvasInfo() SlideCanvasInfo {
	return SlideCanvasInfo{
		BoxFgColor:      MustParseRgbColorHex("#ffffff"),
		BoxBgColor:      MustParseRgbColorHex("#4688f1"),
		BoxHeight:       25.0,
		BoxMarginBottom: 5.0,
		Canvas: CanvasFloat64{
			MinX: 150.0,
			MinY: 70.0,
			MaxX: 700.0,
			MaxY: 500.0,
		},
	}
}

type CanvasFloat64 struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

func (c64 *CanvasFloat64) ThisX(this, min, max float64) (float64, error) {
	if min > max {
		return 0.0, fmt.Errorf("min [%v] is larger than max [%v]", min, max)
	} else if this < min || this > max {
		return 0.0, fmt.Errorf("this [%v] is not within min,max [%v, %v]", this, min, max)
	}
	diff := max - min
	plus := this - min
	pct := plus / diff
	diffCan := c64.MaxX - c64.MinX
	thisPlus := pct * diffCan
	thisX := c64.MinX + thisPlus
	return thisX, nil
}

type Location struct {
	SrcAllMinX int64
	SrcAllMaxX int64
	SrcAllWdtX int64
	SrcBoxMinX int64
	SrcBoxMaxX int64
	SrcBoxWdtX int64
	SrcPctWdtX float64
	OutAllMinX float64
	OutAllMaxX float64
	OutAllWdtX float64
	OutBoxMinX float64
	OutBoxMaxX float64
	OutBoxWdtX float64
	BoxOutPctX float64
}

func GoogleSlideDrawRoadmap(pageID string, srcCan roadmap.Canvas, outCan SlideCanvasInfo) ([]*slides.Request, error) {
	requests := []*slides.Request{}
	err := srcCan.InflateItems()
	if err != nil {
		return requests, err
	}
	srcCan.BuildRows()

	idx := 0
	rowYWatermark := outCan.Canvas.MinY

	for _, row := range srcCan.Rows {
		for _, el := range row {
			// fmtutil.PrintJSON(el)
			srcBoxWdtX := el.Max - el.Min
			srcAllWdtX := srcCan.MaxX - srcCan.MinX
			srcBoxMinX := el.Min
			srcBoxMaxX := el.Max
			srcPctWdtX := float64(srcBoxWdtX) / float64(srcAllWdtX)

			srcAllMinX := srcCan.MinX
			outAllWdtX := outCan.Canvas.MaxX - outCan.Canvas.MinX
			outBoxWdtX := srcPctWdtX * outAllWdtX

			boxOutPctX := float64(srcAllWdtX) / outAllWdtX

			outAllMinX := outCan.Canvas.MinX
			// fmt.Printf("%v\n", srcBoxMinX-srcAllMinX)
			// fmt.Printf("%v\n", el.Min-srcCan.MinX)
			outBoxMinX := (float64(srcBoxMinX-srcAllMinX) / boxOutPctX) + outAllMinX
			outBoxMaxX := (float64(srcBoxMaxX-srcAllMinX) / boxOutPctX) + outAllMinX

			loc := Location{
				SrcAllMinX: srcCan.MinX,
				SrcAllMaxX: srcCan.MaxX,
				SrcAllWdtX: srcCan.MaxX - srcCan.MinX,
				SrcBoxMinX: el.Min,
				SrcBoxMaxX: el.Max,
				SrcBoxWdtX: srcBoxWdtX,
				SrcPctWdtX: srcPctWdtX,
				OutAllMinX: outCan.Canvas.MinX,
				OutAllMaxX: outCan.Canvas.MaxX,
				OutAllWdtX: outCan.Canvas.MaxX - outCan.Canvas.MinX,
				OutBoxMinX: outBoxMinX,
				OutBoxMaxX: outBoxMaxX,
				OutBoxWdtX: outBoxWdtX,
				BoxOutPctX: boxOutPctX,
			}

			// fmtutil.PrintJSON(loc)
			if loc.OutBoxMaxX > loc.OutAllMaxX {
				panic("C")
			} else if loc.OutBoxMinX < loc.OutAllMinX {
				panic("D")
			}
			//panic("Z")
			elementID := fmt.Sprintf("AutoBox%03d", idx)
			requests = append(requests, TextBoxRequestsSimple(
				pageID, elementID, el.NameShort, outCan.BoxFgColor, outCan.BoxBgColor,
				loc.OutBoxWdtX, outCan.BoxHeight, loc.OutBoxMinX, rowYWatermark)...)
			idx++
			//break
			/*
				if i > 3 {
					break
				}
			*/
		}
		rowYWatermark += outCan.BoxHeight + outCan.BoxMarginBottom
	}
	return requests, nil
}
