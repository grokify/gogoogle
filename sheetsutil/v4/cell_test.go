package sheetsutil

import (
	"testing"

	"google.golang.org/api/sheets/v4"
)

func TestParseCellValue(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		formatted string
		wantType  CellType
	}{
		{
			name:      "nil value",
			value:     nil,
			formatted: "",
			wantType:  CellTypeEmpty,
		},
		{
			name:      "empty string",
			value:     "",
			formatted: "",
			wantType:  CellTypeEmpty,
		},
		{
			name:      "string value",
			value:     "Hello",
			formatted: "Hello",
			wantType:  CellTypeString,
		},
		{
			name:      "float64 value",
			value:     float64(123.45),
			formatted: "$123.45",
			wantType:  CellTypeNumber,
		},
		{
			name:      "int value",
			value:     42,
			formatted: "42",
			wantType:  CellTypeNumber,
		},
		{
			name:      "int64 value",
			value:     int64(100),
			formatted: "100",
			wantType:  CellTypeNumber,
		},
		{
			name:      "bool true",
			value:     true,
			formatted: "TRUE",
			wantType:  CellTypeBool,
		},
		{
			name:      "bool false",
			value:     false,
			formatted: "FALSE",
			wantType:  CellTypeBool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCellValue(tt.value, tt.formatted)
			if got.Type != tt.wantType {
				t.Errorf("ParseCellValue().Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.FormattedValue != tt.formatted {
				t.Errorf("ParseCellValue().FormattedValue = %v, want %v", got.FormattedValue, tt.formatted)
			}
		})
	}
}

func TestParseCellValueNumberDetails(t *testing.T) {
	cv := ParseCellValue(float64(123.45), "$123.45")

	if cv.Type != CellTypeNumber {
		t.Errorf("Type = %v, want %v", cv.Type, CellTypeNumber)
	}
	if cv.NumberValue == nil {
		t.Error("NumberValue should not be nil")
	} else if *cv.NumberValue != 123.45 {
		t.Errorf("NumberValue = %v, want %v", *cv.NumberValue, 123.45)
	}
	if cv.StringValue != nil {
		t.Error("StringValue should be nil for number")
	}
}

func TestParseCellValueStringDetails(t *testing.T) {
	cv := ParseCellValue("Hello World", "Hello World")

	if cv.Type != CellTypeString {
		t.Errorf("Type = %v, want %v", cv.Type, CellTypeString)
	}
	if cv.StringValue == nil {
		t.Error("StringValue should not be nil")
	} else if *cv.StringValue != "Hello World" {
		t.Errorf("StringValue = %v, want %v", *cv.StringValue, "Hello World")
	}
	if cv.NumberValue != nil {
		t.Error("NumberValue should be nil for string")
	}
}

func TestParseCellValueBoolDetails(t *testing.T) {
	cv := ParseCellValue(true, "TRUE")

	if cv.Type != CellTypeBool {
		t.Errorf("Type = %v, want %v", cv.Type, CellTypeBool)
	}
	if cv.BoolValue == nil {
		t.Error("BoolValue should not be nil")
	} else if !*cv.BoolValue {
		t.Error("BoolValue should be true")
	}
}

func TestParseTypedCellValueWithPreferInt64(t *testing.T) {
	opts := ValueParseOptions{PreferInt64: true}
	tcv := ParseTypedCellValue(float64(42), "42", opts)

	if tcv.Type != CellTypeNumber {
		t.Errorf("Type = %v, want %v", tcv.Type, CellTypeNumber)
	}
	if tcv.Int64 == nil {
		t.Error("Int64 should not be nil for whole number with PreferInt64=true")
	} else if *tcv.Int64 != 42 {
		t.Errorf("Int64 = %v, want %v", *tcv.Int64, 42)
	}
}

func TestParseTypedCellValueFractional(t *testing.T) {
	opts := ValueParseOptions{PreferInt64: true}
	tcv := ParseTypedCellValue(float64(42.5), "42.5", opts)

	if tcv.Type != CellTypeNumber {
		t.Errorf("Type = %v, want %v", tcv.Type, CellTypeNumber)
	}
	if tcv.Int64 != nil {
		t.Error("Int64 should be nil for fractional number")
	}
}

func TestParseValueRange(t *testing.T) {
	vr := &sheets.ValueRange{
		Values: [][]any{
			{"Name", "Age", "Active"},
			{"Alice", float64(30), true},
			{"Bob", float64(25), false},
		},
	}

	grid := ParseValueRange(vr)

	if len(grid) != 3 {
		t.Fatalf("grid length = %d, want 3", len(grid))
	}
	if len(grid[0]) != 3 {
		t.Fatalf("grid[0] length = %d, want 3", len(grid[0]))
	}

	// Check header row
	if grid[0][0].Type != CellTypeString {
		t.Errorf("grid[0][0].Type = %v, want %v", grid[0][0].Type, CellTypeString)
	}
	if grid[0][0].FormattedValue != "Name" {
		t.Errorf("grid[0][0].FormattedValue = %v, want %v", grid[0][0].FormattedValue, "Name")
	}

	// Check data row
	if grid[1][1].Type != CellTypeNumber {
		t.Errorf("grid[1][1].Type = %v, want %v", grid[1][1].Type, CellTypeNumber)
	}
	if grid[1][2].Type != CellTypeBool {
		t.Errorf("grid[1][2].Type = %v, want %v", grid[1][2].Type, CellTypeBool)
	}
}

func TestParseValueRangeNil(t *testing.T) {
	grid := ParseValueRange(nil)
	if grid != nil {
		t.Error("ParseValueRange(nil) should return nil")
	}
}

func TestParseValueRangeEmpty(t *testing.T) {
	vr := &sheets.ValueRange{Values: [][]any{}}
	grid := ParseValueRange(vr)
	if grid != nil {
		t.Error("ParseValueRange with empty values should return nil")
	}
}

func TestExtractFormattedValues(t *testing.T) {
	vr := &sheets.ValueRange{
		Values: [][]any{
			{"Name", "Amount"},
			{"Alice", float64(1234.56)},
			{"Bob", true},
		},
	}

	result := ExtractFormattedValues(vr)

	if len(result) != 3 {
		t.Fatalf("result length = %d, want 3", len(result))
	}

	// Check string preservation
	if result[0][0] != "Name" {
		t.Errorf("result[0][0] = %v, want Name", result[0][0])
	}

	// Check bool formatting
	if result[2][1] != "TRUE" {
		t.Errorf("result[2][1] = %v, want TRUE", result[2][1])
	}
}

func TestExtractFormattedValuesNil(t *testing.T) {
	result := ExtractFormattedValues(nil)
	if result != nil {
		t.Error("ExtractFormattedValues(nil) should return nil")
	}
}

func TestExtractRawValues(t *testing.T) {
	vr := &sheets.ValueRange{
		Values: [][]any{
			{"Name", "Amount"},
			{"Alice", float64(1234)},
			{"Bob", false},
		},
	}

	result := ExtractRawValues(vr)

	if len(result) != 3 {
		t.Fatalf("result length = %d, want 3", len(result))
	}

	// Check number formatting (should be raw, no commas)
	if result[1][1] != "1234" {
		t.Errorf("result[1][1] = %v, want 1234", result[1][1])
	}

	// Check bool formatting
	if result[2][1] != "FALSE" {
		t.Errorf("result[2][1] = %v, want FALSE", result[2][1])
	}
}

func TestGridToStringSlice(t *testing.T) {
	grid := SimpleGrid{
		{
			{Type: CellTypeString, FormattedValue: "Name"},
			{Type: CellTypeString, FormattedValue: "Value"},
		},
		{
			{Type: CellTypeString, FormattedValue: "Test"},
			{Type: CellTypeNumber, FormattedValue: "$100.00"},
		},
	}

	result := GridToStringSlice(grid)

	if len(result) != 2 {
		t.Fatalf("result length = %d, want 2", len(result))
	}
	if result[0][0] != "Name" {
		t.Errorf("result[0][0] = %v, want Name", result[0][0])
	}
	if result[1][1] != "$100.00" {
		t.Errorf("result[1][1] = %v, want $100.00", result[1][1])
	}
}

func TestGridToStringSliceNil(t *testing.T) {
	result := GridToStringSlice(nil)
	if result != nil {
		t.Error("GridToStringSlice(nil) should return nil")
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{123, "123"},
		{-456, "-456"},
		{1.5, "1.5"},
		{3.14159, "3.14159"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatFloat(tt.input)
			if got != tt.want {
				t.Errorf("formatFloat(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
