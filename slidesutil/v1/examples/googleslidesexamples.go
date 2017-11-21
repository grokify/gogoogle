package slidesutilexamples

import (
	"google.golang.org/api/slides/v1"
)

func ExampleRequests() []*slides.Request {
	// https://developers.google.com/slides/samples/tables

	tableId := "tableId"
	pageId := "pageId"

	return []*slides.Request{
		// Create a table
		{
			CreateTable: &slides.CreateTableRequest{
				ObjectId: tableId,
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: pageId,
				},
				Rows:    8,
				Columns: 5,
			},
		},
		{
			// Delete table rows or columns
			DeleteTableRow: &slides.DeleteTableRowRequest{
				TableObjectId: tableId,
				CellLocation: &slides.TableCellLocation{
					RowIndex: 5,
				},
			},
		},
		{
			DeleteTableColumn: &slides.DeleteTableColumnRequest{
				TableObjectId: tableId,
				CellLocation: &slides.TableCellLocation{
					ColumnIndex: 3,
				},
			},
		},
		{
			// Edit table data
			DeleteText: &slides.DeleteTextRequest{
				ObjectId: tableId,
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
				ObjectId: "tableId",
				CellLocation: &slides.TableCellLocation{
					ColumnIndex: 4,
					RowIndex:    2,
				},
				Text:           "Kangaroo",
				InsertionIndex: 0,
			},
		},
		{
			// Format a table header row
			UpdateTableCellProperties: &slides.UpdateTableCellPropertiesRequest{
				ObjectId: tableId,
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
				ObjectId: tableId,
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
		{
			// Insert table rows or columns
			InsertTableRows: &slides.InsertTableRowsRequest{
				TableObjectId: tableId,
				CellLocation: &slides.TableCellLocation{
					RowIndex: 5,
				},
				InsertBelow: true,
				Number:      3,
			},
		},
		{
			InsertTableColumns: &slides.InsertTableColumnsRequest{
				TableObjectId: tableId,
				CellLocation: &slides.TableCellLocation{
					ColumnIndex: 3,
				},
				InsertRight: false,
				Number:      2,
			},
		},
	}
}
