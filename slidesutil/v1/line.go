package slidesutil

import (
	"google.golang.org/api/slides/v1"
)

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
			c, err := ParseRgbColorHex(info.ColorHex)
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
