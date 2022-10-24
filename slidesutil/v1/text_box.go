package slidesutil

import (
	"strings"

	slides "google.golang.org/api/slides/v1"
)

type CreateShapeTextBoxRequestInfo struct {
	PageId             string
	ObjectId           string
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

func (info *CreateShapeTextBoxRequestInfo) NeedsUpdateTextStyle() bool {
	return info.ForegroundColorRgb != nil ||
		len(info.ForegroundColorHex) > 0 ||
		info.FontSize > 0.0 ||
		info.FontBold ||
		info.FontItalic
}

func (info *CreateShapeTextBoxRequestInfo) NeedsUpdateShapeProperties() bool {
	return info.BackgroundColorRgb != nil ||
		len(info.BackgroundColorHex) > 0 || len(info.URL) > 0
}

func (info *CreateShapeTextBoxRequestInfo) Requests() ([]*slides.Request, error) {
	requests := []*slides.Request{
		{
			CreateShape: &slides.CreateShapeRequest{
				ObjectId:  info.ObjectId,
				ShapeType: "TEXT_BOX",
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: info.PageId,
					Size: &slides.Size{
						Width:  &slides.Dimension{Magnitude: info.Width, Unit: info.DimensionUnit},
						Height: &slides.Dimension{Magnitude: info.Height, Unit: info.DimensionUnit},
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: info.LocationX,
						TranslateY: info.LocationY,
						Unit:       info.LocationUnit,
					},
				},
			},
		},
	}
	if len(info.Text) > 0 {
		requests = append(requests, &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       info.ObjectId,
				InsertionIndex: 0,
				Text:           info.Text,
			},
		})
	}

	if info.ForegroundColorRgb != nil ||
		len(info.ForegroundColorHex) > 0 ||
		info.FontSize > 0.0 ||
		info.FontBold || info.FontItalic {
		req := &slides.Request{
			UpdateTextStyle: &slides.UpdateTextStyleRequest{
				ObjectId: info.ObjectId,
				Fields:   "*",
				Style:    &slides.TextStyle{},
			},
		}

		if info.ForegroundColorRgb != nil {
			req.UpdateTextStyle.Style.ForegroundColor = &slides.OptionalColor{
				OpaqueColor: &slides.OpaqueColor{
					RgbColor: info.ForegroundColorRgb,
				},
			}
		} else if len(info.ForegroundColorHex) > 0 {
			c, err := ParseRgbColorHex(info.ForegroundColorHex)
			if err != nil {
				return requests, err
			}
			req.UpdateTextStyle.Style.ForegroundColor = &slides.OptionalColor{
				OpaqueColor: &slides.OpaqueColor{
					RgbColor: c,
				},
			}
		}
		if info.FontSize > 0.0 {
			req.UpdateTextStyle.Style.FontSize = &slides.Dimension{
				Magnitude: info.FontSize,
				Unit:      info.FontSizeUnit,
			}
		}
		if info.FontBold {
			req.UpdateTextStyle.Style.Bold = true
		}
		if info.FontItalic {
			req.UpdateTextStyle.Style.Italic = true
		}
		requests = append(requests, req)
	}

	if len(info.ParagraphAlignment) > 0 {
		req := &slides.Request{
			UpdateParagraphStyle: &slides.UpdateParagraphStyleRequest{
				ObjectId: info.ObjectId,
				Fields:   "*",
				Style: &slides.ParagraphStyle{
					Alignment: info.ParagraphAlignment,
				},
			},
		}
		requests = append(requests, req)
	}

	if info.NeedsUpdateShapeProperties() {
		shapeUtil := NewUpdateShapePropertiesRequestUtil(info.ObjectId)

		if info.BackgroundColorRgb != nil {
			shapeUtil.AddBackgroundSolidFill(info.BackgroundColorRgb)
		} else if len(info.BackgroundColorHex) > 0 {
			c, err := ParseRgbColorHex(info.BackgroundColorHex)
			if err != nil {
				return requests, err
			}
			shapeUtil.AddBackgroundSolidFill(c)
		}
		if len(info.URL) > 0 {
			shapeUtil.AddLink(&slides.Link{Url: info.URL})
		}
		requests = append(requests, shapeUtil.Request())
	}

	if 1 == 0 {
		if info.BackgroundColorRgb != nil {
			requests = append(requests,
				ShapePropertiesBackgroundFillSimple(info.ObjectId, info.BackgroundColorRgb))
		} else if len(info.BackgroundColorHex) > 0 {
			c, err := ParseRgbColorHex(info.BackgroundColorHex)
			if err != nil {
				return requests, err
			}
			requests = append(requests,
				ShapePropertiesBackgroundFillSimple(info.ObjectId, c))
		}
	}
	/*
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
	},*/

	return requests, nil
}

type UpdateShapePropertiesRequestUtil struct {
	Fields                       []string
	UpdateShapePropertiesRequest *slides.UpdateShapePropertiesRequest
}

func NewUpdateShapePropertiesRequestUtil(objectId string) UpdateShapePropertiesRequestUtil {
	return UpdateShapePropertiesRequestUtil{
		Fields: []string{},
		UpdateShapePropertiesRequest: &slides.UpdateShapePropertiesRequest{
			ObjectId: objectId,
		},
	}
}

func (util *UpdateShapePropertiesRequestUtil) AddBackgroundSolidFill(rgbColor *slides.RgbColor) {
	if util.UpdateShapePropertiesRequest.ShapeProperties == nil {
		util.UpdateShapePropertiesRequest.ShapeProperties = &slides.ShapeProperties{}
	}
	util.UpdateShapePropertiesRequest.ShapeProperties.ShapeBackgroundFill = &slides.ShapeBackgroundFill{
		SolidFill: &slides.SolidFill{
			Color: &slides.OpaqueColor{
				RgbColor: rgbColor,
			},
		},
	}
	util.Fields = append(util.Fields, "shapeBackgroundFill.solidFill.color")
}

func (util *UpdateShapePropertiesRequestUtil) AddLink(link *slides.Link) {
	if util.UpdateShapePropertiesRequest.ShapeProperties == nil {
		util.UpdateShapePropertiesRequest.ShapeProperties = &slides.ShapeProperties{}
	}
	util.UpdateShapePropertiesRequest.ShapeProperties.Link = link
	util.Fields = append(util.Fields, "link")
}

func (util *UpdateShapePropertiesRequestUtil) Request() *slides.Request {
	if len(util.Fields) > 0 {
		util.UpdateShapePropertiesRequest.Fields = strings.Join(util.Fields, ",")
	}
	return &slides.Request{
		UpdateShapeProperties: util.UpdateShapePropertiesRequest,
	}
}

// ShapePropertiesBackgroundFillSimple is a simple shape properties request
// creator. For more complex uses, use UpdateShapePropertiesRequestUtil
func ShapePropertiesBackgroundFillSimple(objectId string, rgbColor *slides.RgbColor) *slides.Request {
	return &slides.Request{
		UpdateShapeProperties: &slides.UpdateShapePropertiesRequest{
			ObjectId: objectId,
			Fields:   "shapeBackgroundFill.solidFill.color",
			ShapeProperties: &slides.ShapeProperties{
				ShapeBackgroundFill: &slides.ShapeBackgroundFill{
					SolidFill: &slides.SolidFill{
						Color: &slides.OpaqueColor{
							RgbColor: rgbColor,
						},
					},
				},
			},
		},
	}
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
