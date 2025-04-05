package slidesutilexamples

import (
	"google.golang.org/api/slides/v1"

	su "github.com/grokify/gogoogle/slidesutil/v1"
)

func LineExampleRequests(pageID string) []*slides.Request {
	lineID := "lineId"

	return []*slides.Request{
		{
			CreateLine: &slides.CreateLineRequest{
				ObjectId:     lineID,
				LineCategory: "STRAIGHT",
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: pageID,
					Size: &slides.Size{
						Height: &slides.Dimension{Magnitude: 500.0, Unit: "PT"},
						Width:  &slides.Dimension{Magnitude: 1.0, Unit: "PT"},
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: 350.0,
						TranslateY: 100.0,
						Unit:       "PT",
					},
				},
			},
		},
		{
			UpdateLineProperties: &slides.UpdateLinePropertiesRequest{
				ObjectId: lineID,
				Fields:   "*",
				LineProperties: &slides.LineProperties{
					DashStyle: "DASH",
					LineFill: &slides.LineFill{
						SolidFill: &slides.SolidFill{
							Color: &slides.OpaqueColor{
								RgbColor: su.MustParseRgbColorHex("#ff8800"),
							},
						},
					},
					Weight: &slides.Dimension{
						Magnitude: 1.0,
						Unit:      "PT",
					},
				},
			},
		},
	}
}

// ExampleRequests provides Go code for JSON examples provided by Google
// at https://developers.google.com/slides/samples/tables
func TableExampleRequests() []*slides.Request {
	tableID := "tableId"
	pageID := "pageId"

	return []*slides.Request{
		// Create a table
		{
			CreateTable: &slides.CreateTableRequest{
				ObjectId: tableID,
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: pageID,
				},
				Rows:    8,
				Columns: 5,
			},
		},
		// Delete table rows or columns
		{
			DeleteTableRow: &slides.DeleteTableRowRequest{
				TableObjectId: tableID,
				CellLocation: &slides.TableCellLocation{
					RowIndex: 5,
				},
			},
		},
		{
			DeleteTableColumn: &slides.DeleteTableColumnRequest{
				TableObjectId: tableID,
				CellLocation: &slides.TableCellLocation{
					ColumnIndex: 3,
				},
			},
		},
		// Edit table data
		{
			DeleteText: &slides.DeleteTextRequest{
				ObjectId: tableID,
				CellLocation: &slides.TableCellLocation{
					ColumnIndex: 4,
					RowIndex:    2,
				},
				TextRange: &slides.Range{
					Type: "ALL",
				},
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId: tableID,
				CellLocation: &slides.TableCellLocation{
					ColumnIndex: 4,
					RowIndex:    2,
				},
				Text:           "Kangaroo",
				InsertionIndex: 0,
			},
		},
		// Format a table header row
		{
			UpdateTableCellProperties: &slides.UpdateTableCellPropertiesRequest{
				ObjectId: tableID,
				TableRange: &slides.TableRange{
					Location: &slides.TableCellLocation{
						RowIndex:    0,
						ColumnIndex: 0,
					},
					RowSpan:    1,
					ColumnSpan: 1,
				},
				TableCellProperties: &slides.TableCellProperties{
					TableCellBackgroundFill: &slides.TableCellBackgroundFill{
						SolidFill: &slides.SolidFill{
							Color: &slides.OpaqueColor{
								RgbColor: &slides.RgbColor{
									Red:   0.0,
									Green: 0.0,
									Blue:  0.0,
								},
							},
						},
					},
				},
				Fields: "tableCellBackgroundFill.solidFill.color",
			},
		},
		{
			UpdateTextStyle: &slides.UpdateTextStyleRequest{
				ObjectId: tableID,
				CellLocation: &slides.TableCellLocation{
					RowIndex:    0,
					ColumnIndex: 0,
				},
				Style: &slides.TextStyle{
					ForegroundColor: &slides.OptionalColor{
						OpaqueColor: &slides.OpaqueColor{
							RgbColor: &slides.RgbColor{
								Red:   1.0,
								Green: 1.0,
								Blue:  1.0,
							},
						},
					},
					Bold:       true,
					FontFamily: "Cambria",
					FontSize: &slides.Dimension{
						Magnitude: 18,
						Unit:      "PT",
					},
				},
				TextRange: &slides.Range{
					Type: "ALL",
				},
				Fields: "foregroundColor,bold,fontFamily,fontSize",
			},
		},
		// Insert table rows or columns
		{
			InsertTableRows: &slides.InsertTableRowsRequest{
				TableObjectId: tableID,
				CellLocation: &slides.TableCellLocation{
					RowIndex: 5,
				},
				InsertBelow: true,
				Number:      3,
			},
		},
		{
			InsertTableColumns: &slides.InsertTableColumnsRequest{
				TableObjectId: tableID,
				CellLocation: &slides.TableCellLocation{
					ColumnIndex: 3,
				},
				InsertRight: false,
				Number:      2,
			},
		},
	}
}
