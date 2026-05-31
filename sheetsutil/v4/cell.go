package sheetsutil

import (
	"math"
	"strconv"
	"time"

	"google.golang.org/api/sheets/v4"
)

// CellType represents the type of data in a cell.
type CellType string

const (
	CellTypeEmpty    CellType = "empty"
	CellTypeString   CellType = "string"
	CellTypeNumber   CellType = "number"
	CellTypeBool     CellType = "boolean"
	CellTypeDate     CellType = "date"
	CellTypeTime     CellType = "time"
	CellTypeDateTime CellType = "datetime"
	CellTypeDuration CellType = "duration"
	CellTypeError    CellType = "error"
)

// CellValue represents a simple, JSON-friendly cell value.
type CellValue struct {
	Type           CellType `json:"type"`
	FormattedValue string   `json:"formatted_value"`
	StringValue    *string  `json:"string_value,omitempty"`
	NumberValue    *float64 `json:"number_value,omitempty"`
	BoolValue      *bool    `json:"bool_value,omitempty"`
	ErrorValue     *string  `json:"error_value,omitempty"`
}

// TypedCellValue extends CellValue with Go-native types for dates and durations.
type TypedCellValue struct {
	CellValue
	Time     *time.Time     `json:"-"`
	Duration *time.Duration `json:"-"`
	Int64    *int64         `json:"-"`
}

// ValueParseOptions configures how values are parsed.
type ValueParseOptions struct {
	Timezone    *time.Location
	DateFormats []string
	PreferInt64 bool
}

// DefaultValueParseOptions returns default parsing options.
func DefaultValueParseOptions() ValueParseOptions {
	return ValueParseOptions{
		Timezone:    time.UTC,
		PreferInt64: false,
	}
}

// SimpleGrid is a 2D array of CellValue for JSON serialization.
type SimpleGrid = [][]CellValue

// TypedGrid is a 2D array of TypedCellValue for Go-native processing.
type TypedGrid = [][]TypedCellValue

// ParseCellValue converts a raw API value to a CellValue.
// The value parameter is the raw value from the Sheets API (typically any type).
// The formattedValue is the user-visible formatted string.
func ParseCellValue(value any, formattedValue string) CellValue {
	cv := CellValue{
		FormattedValue: formattedValue,
	}

	if value == nil {
		cv.Type = CellTypeEmpty
		return cv
	}

	switch v := value.(type) {
	case string:
		if v == "" && formattedValue == "" {
			cv.Type = CellTypeEmpty
		} else {
			cv.Type = CellTypeString
			cv.StringValue = &v
		}
	case float64:
		cv.Type = CellTypeNumber
		cv.NumberValue = &v
	case int:
		cv.Type = CellTypeNumber
		f := float64(v)
		cv.NumberValue = &f
	case int64:
		cv.Type = CellTypeNumber
		f := float64(v)
		cv.NumberValue = &f
	case bool:
		cv.Type = CellTypeBool
		cv.BoolValue = &v
	default:
		// Treat unknown types as strings
		s := formattedValue
		cv.Type = CellTypeString
		cv.StringValue = &s
	}

	return cv
}

// ParseTypedCellValue converts a raw API value to a TypedCellValue with richer type information.
func ParseTypedCellValue(value any, formattedValue string, opts ValueParseOptions) TypedCellValue {
	cv := ParseCellValue(value, formattedValue)
	tcv := TypedCellValue{CellValue: cv}

	if cv.NumberValue != nil && opts.PreferInt64 {
		f := *cv.NumberValue
		if f == math.Trunc(f) && f >= math.MinInt64 && f <= math.MaxInt64 {
			i := int64(f)
			tcv.Int64 = &i
		}
	}

	return tcv
}

// ParseValueRange converts a Sheets API ValueRange to a SimpleGrid.
func ParseValueRange(vr *sheets.ValueRange) SimpleGrid {
	if vr == nil || len(vr.Values) == 0 {
		return nil
	}

	grid := make(SimpleGrid, len(vr.Values))
	for i, row := range vr.Values {
		grid[i] = make([]CellValue, len(row))
		for j, val := range row {
			// In ValueRenderOption FORMATTED_VALUE, values are strings
			// In ValueRenderOption UNFORMATTED_VALUE, values retain their types
			formatted := ""
			if s, ok := val.(string); ok {
				formatted = s
			}
			grid[i][j] = ParseCellValue(val, formatted)
		}
	}
	return grid
}

// ParseTypedValueRange converts a Sheets API ValueRange to a TypedGrid.
func ParseTypedValueRange(vr *sheets.ValueRange, opts ValueParseOptions) TypedGrid {
	if vr == nil || len(vr.Values) == 0 {
		return nil
	}

	grid := make(TypedGrid, len(vr.Values))
	for i, row := range vr.Values {
		grid[i] = make([]TypedCellValue, len(row))
		for j, val := range row {
			formatted := ""
			if s, ok := val.(string); ok {
				formatted = s
			}
			grid[i][j] = ParseTypedCellValue(val, formatted, opts)
		}
	}
	return grid
}

// ExtractFormattedValues extracts just the formatted value strings from a ValueRange.
// This is useful when you want a simple [][]string representation.
func ExtractFormattedValues(vr *sheets.ValueRange) [][]string {
	if vr == nil || len(vr.Values) == 0 {
		return nil
	}

	result := make([][]string, len(vr.Values))
	for i, row := range vr.Values {
		result[i] = make([]string, len(row))
		for j, val := range row {
			switch v := val.(type) {
			case string:
				result[i][j] = v
			case float64:
				// For numbers, use the formatted value if available
				result[i][j] = formatFloat(v)
			case bool:
				if v {
					result[i][j] = "TRUE"
				} else {
					result[i][j] = "FALSE"
				}
			default:
				result[i][j] = ""
			}
		}
	}
	return result
}

// ExtractRawValues extracts raw values as strings from a ValueRange.
// Numbers are converted without formatting.
func ExtractRawValues(vr *sheets.ValueRange) [][]string {
	if vr == nil || len(vr.Values) == 0 {
		return nil
	}

	result := make([][]string, len(vr.Values))
	for i, row := range vr.Values {
		result[i] = make([]string, len(row))
		for j, val := range row {
			switch v := val.(type) {
			case string:
				result[i][j] = v
			case float64:
				result[i][j] = formatFloat(v)
			case bool:
				if v {
					result[i][j] = "TRUE"
				} else {
					result[i][j] = "FALSE"
				}
			default:
				result[i][j] = ""
			}
		}
	}
	return result
}

// formatFloat formats a float64, removing trailing zeros.
// Uses strconv.FormatFloat with 'g' format for clean output.
func formatFloat(f float64) string {
	// Handle special cases
	if math.IsNaN(f) {
		return "NaN"
	}
	if math.IsInf(f, 1) {
		return "+Inf"
	}
	if math.IsInf(f, -1) {
		return "-Inf"
	}

	// For whole numbers, format as integer
	if f == math.Trunc(f) && f >= -9007199254740992 && f <= 9007199254740992 {
		return strconv.FormatInt(int64(f), 10)
	}

	// Use 'g' format which removes trailing zeros and uses scientific notation when appropriate
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// GridToStringSlice converts a SimpleGrid to [][]string using formatted values.
func GridToStringSlice(grid SimpleGrid) [][]string {
	if grid == nil {
		return nil
	}
	result := make([][]string, len(grid))
	for i, row := range grid {
		result[i] = make([]string, len(row))
		for j, cell := range row {
			result[i][j] = cell.FormattedValue
		}
	}
	return result
}
