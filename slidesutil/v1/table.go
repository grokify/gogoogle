package slidesutil

import (
	"strings"

	"google.golang.org/api/slides/v1"
)

func IsEven(i int) bool {
	return i%2 == 0
}

func IsOdd(i int) bool {
	return !IsEven(i)
}

func AlternateRowBgColor(objectID string, rowCount, columnCount int64, evenColorHex, oddColorHex string) ([]*slides.Request, error) {
	reqs := []*slides.Request{}

	colorizeEven := false
	colorizeOdd := false
	evenColor := &slides.RgbColor{}
	oddColor := &slides.RgbColor{}

	if len(evenColorHex) > 0 {
		c, err := ParseRgbColorHex(evenColorHex)
		if err != nil {
			return reqs, err
		}
		evenColor = c
		colorizeEven = true
	}
	if len(oddColorHex) > 0 {
		c, err := ParseRgbColorHex(oddColorHex)
		if err != nil {
			return reqs, err
		}
		oddColor = c
		colorizeOdd = true
	}

	for i := 0; i < int(rowCount); i++ {
		if IsEven(i) && colorizeEven {
			reqs = append(reqs,
				&slides.Request{
					UpdateTableCellProperties: UpdateTableCellPropertiesRequestTableCellBackgroundFill(objectID, int64(i), columnCount, evenColor),
				},
			)
		} else if IsOdd(i) && colorizeOdd {
			reqs = append(reqs,
				&slides.Request{
					UpdateTableCellProperties: UpdateTableCellPropertiesRequestTableCellBackgroundFill(objectID, int64(i), columnCount, oddColor),
				},
			)
		}
	}
	return reqs, nil
}

func UpdateTableCellPropertiesRequestTableCellBackgroundFill(objectID string, rowIndex, columnSpan int64, bgColor *slides.RgbColor) *slides.UpdateTableCellPropertiesRequest {
	return &slides.UpdateTableCellPropertiesRequest{
		ObjectId: objectID,
		Fields:   "*",
		TableRange: &slides.TableRange{
			Location: &slides.TableCellLocation{
				ColumnIndex: 0,
				RowIndex:    rowIndex,
			},
			ColumnSpan: columnSpan,
			RowSpan:    1,
		},
		TableCellProperties: &slides.TableCellProperties{
			TableCellBackgroundFill: &slides.TableCellBackgroundFill{
				SolidFill: &slides.SolidFill{
					Color: &slides.OpaqueColor{
						RgbColor: bgColor,
					},
				},
			},
		},
	}
}

// https://developers.google.com/slides/samples/tables#format_a_table_header_row

type UpdateTextStyle struct {
	ObjectID           string
	RowIndex           int64
	ColumnIndex        int64
	ForegroundColorHex string
	Bold               bool
	FontFamily         string
	FontSizeMagnitude  float64
	FontSizeUnit       string
	TextRangeType      string
	Fields             string
}

func (item *UpdateTextStyle) RequestsColumnSpan(columnSpan int64) ([]*slides.Request, error) {
	reqs := []*slides.Request{}
	for i := 0; i < int(columnSpan); i++ {
		item.ColumnIndex = int64(i)
		req, err := item.Request()
		if err != nil {
			return reqs, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func (item *UpdateTextStyle) Request() (*slides.Request, error) {
	fields := []string{"bold"}
	req := &slides.UpdateTextStyleRequest{
		ObjectId: item.ObjectID,
		CellLocation: &slides.TableCellLocation{
			ColumnIndex: item.ColumnIndex,
			RowIndex:    item.RowIndex,
		},
		Style: &slides.TextStyle{
			Bold: item.Bold,
		},
		TextRange: &slides.Range{
			Type: item.TextRangeType,
		},
		Fields: item.Fields,
	}
	if len(item.ForegroundColorHex) > 0 {
		c, err := ParseRgbColorHex(item.ForegroundColorHex)
		if err != nil {
			return nil, err
		}
		req.Style.ForegroundColor = &slides.OptionalColor{
			OpaqueColor: &slides.OpaqueColor{
				RgbColor: c,
			},
		}
		fields = append(fields, "foregroundColor")
	}
	if item.FontSizeMagnitude > 0 && len(item.FontSizeUnit) > 0 {
		req.Style.FontSize = &slides.Dimension{
			Magnitude: item.FontSizeMagnitude,
			Unit:      item.FontSizeUnit,
		}
		fields = append(fields, "fontSize")
	}
	req.Fields = strings.Join(fields, ",")

	return &slides.Request{UpdateTextStyle: req}, nil
}
