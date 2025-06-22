package iwark

import (
	"github.com/Iwark/spreadsheet"
)

type Sheet spreadsheet.Sheet

// SheetUpdateMulti is a convenience function that wraps `spreadsheet.Sheet.Update()`
// and `spreadsheet.Sheet.Synchronize()`.
func SheetUpdateMulti(sheet *spreadsheet.Sheet, cells []Cell, sync bool) error {
	if sheet == nil {
		return ErrSheetCannotBeNil
	}
	for _, cell := range cells {
		if cInt, rInt, err := cell.ColumnRowInts(); err != nil {
			return err
		} else {
			sheet.Update(rInt, cInt, cell.Value)
		}
	}
	if sync && len(cells) > 0 {
		return sheet.Synchronize()
	} else {
		return nil
	}
}
