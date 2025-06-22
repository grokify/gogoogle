package iwark

import (
	"github.com/Iwark/spreadsheet"
	"github.com/grokify/mogo/type/number"
)

type Cells []Cell

func (cs Cells) Values() []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.Value)
	}
	return out
}

type Cell spreadsheet.Cell

func (c Cell) ColumnInt() (int, error) {
	return number.Utoi(c.Column)
}

func (c Cell) ColumnRowInts() (col, row int, err error) {
	if col, err = c.ColumnInt(); err != nil {
		return
	}
	row, err = c.RowInt()
	return
}

func (c Cell) RowInt() (int, error) {
	return number.Utoi(c.Row)
}
