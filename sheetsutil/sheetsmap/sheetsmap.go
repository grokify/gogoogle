package sheetsmap

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Iwark/spreadsheet"
	"github.com/grokify/gotilla/strings/stringsutil"
)

type SheetsMap struct {
	GoogleClient   *http.Client
	Service        *spreadsheet.Service
	sheetsId       string
	Spreadsheet    spreadsheet.Spreadsheet
	Sheet          *spreadsheet.Sheet
	sheetIndex     uint
	KeyColumnIndex uint
	Columns        []Column
	ColumnMapKeyLc map[string]Column
	ItemMap        map[string]Item
}

func (sm *SheetsMap) ColumnsKeys() []string {
	keys := []string{}
	for _, col := range sm.Columns {
		keys = append(keys, col.Value)
	}
	return keys
}

func NewSheetsMap(client *http.Client, spreadsheetId string, sheetIndex uint) (SheetsMap, error) {
	sm := SheetsMap{
		GoogleClient:   client,
		Service:        spreadsheet.NewServiceWithClient(client),
		sheetsId:       spreadsheetId,
		sheetIndex:     sheetIndex,
		Columns:        []Column{},
		ColumnMapKeyLc: map[string]Column{},
		ItemMap:        map[string]Item{},
	}

	spreadsheet, err := sm.Service.FetchSpreadsheet(spreadsheetId)
	if err != nil {
		return sm, err
	}
	sm.Spreadsheet = spreadsheet

	sheet, err := sm.Spreadsheet.SheetByIndex(sheetIndex)
	if err != nil {
		return sm, err
	}
	sm.Sheet = sheet

	return sm, nil
}

type Item struct {
	Key  string
	Row  uint
	Data map[string]string
}

type Column struct {
	Value              string
	Index              uint64
	Enums              []Enum
	AliasLcToCanonical map[string]string
}

func NewColumn() Column {
	return Column{
		Enums:              []Enum{},
		AliasLcToCanonical: map[string]string{}}
}

func (col *Column) AddEnum(enum Enum) {
	col.Enums = append(col.Enums, enum)
	col.AliasLcToCanonical[strings.ToLower(enum.Canonical)] = enum.Canonical
	for _, alias := range enum.Aliases {
		col.AliasLcToCanonical[strings.ToLower(alias)] = enum.Canonical
	}
}

func (col *Column) EnumsCanonical() []string {
	canonicals := []string{}
	for _, enum := range col.Enums {
		canonicals = append(canonicals, enum.Canonical)
	}
	return canonicals
}

type Enum struct {
	Canonical string
	Aliases   []string
}

// ParseColumn
// tshirt size - XS, S, M, L, XL, XXL, XXXL
func ParseColumn(input string) (Column, error) {
	parts := strings.Split(input, " - ")
	col := NewColumn()
	if len(parts) <= 0 {
		return col, fmt.Errorf("Column Format Error for: %v", input)
	} else if len(parts) > 2 {
		return col, fmt.Errorf("Column Format Error for: %v", input)
	}
	col.Value = strings.TrimSpace(parts[0])
	if len(parts) == 2 { // Have enum values
		enums := stringsutil.SplitCondenseSpace(parts[1], ",")
		for _, enumPlus := range enums {
			enumVariations := stringsutil.SplitCondenseSpace(enumPlus, "|")
			if len(enumVariations) > 0 {
				enum := Enum{Canonical: enumVariations[0]}
				if len(enumVariations) > 1 { // Have aliases
					enum.Aliases = enumVariations[1:]
				}
				col.AddEnum(enum)
			}
		}
	}
	return col, nil
}

func (col *Column) ValueToCanonical(val string) (string, error) {
	if len(col.Enums) == 0 {
		return val, nil
	}
	valLc := TrimSpaceToLower(val)

	for tryLc, tryCanonical := range col.AliasLcToCanonical {
		tryLc = TrimSpaceToLower(tryLc)
		if valLc == tryLc {
			return strings.TrimSpace(tryCanonical), nil
		}
	}

	enums := strings.Join(col.EnumsCanonical(), ", ")
	return enums, fmt.Errorf("Column [%v] Value [%v] not valid [%v]", col.Value, val, enums)
}

func (sm *SheetsMap) FullRead() error {
	fmt.Println("FR0")
	err := sm.ReadColumns()
	if err != nil {
		fmt.Println("FR1")
		return err
	}
	fmt.Println("FR2")
	return sm.ReadItems()
}

func (sm *SheetsMap) ReadColumns() error {
	colsMap := map[string]Column{}
	colsArr := []Column{}

	for _, row := range sm.Sheet.Rows {
		for j, cell := range row {
			colValRaw := strings.TrimSpace(cell.Value)
			if len(colValRaw) < 1 {
				break
			}

			col, err := ParseColumn(colValRaw)
			if err != nil {
				return err
			}
			col.Index = uint64(j)
			colKeyParsedLc := strings.ToLower(col.Value)

			colsArr = append(colsArr, col)
			if _, ok := colsMap[colKeyParsedLc]; ok {
				return fmt.Errorf("Duplicate column names for: %v", colValRaw)
			}
			colsMap[colKeyParsedLc] = col
		}
		break
	}

	sm.Columns = colsArr
	sm.ColumnMapKeyLc = colsMap
	return nil
}

func (sm *SheetsMap) ReadItems() error {
	itemMap := map[string]Item{}

	for i, row := range sm.Sheet.Rows {
		if i == 0 {
			continue
		}
		item := Item{
			Row:  uint(i),
			Data: map[string]string{},
		}
		for j, cell := range row {
			if j >= len(sm.Columns) {
				break
			}

			val := cell.Value
			if j == 0 {
				item.Key = val
			}
			col := sm.Columns[j]
			item.Data[col.Value] = val
		}
		if _, ok := itemMap[item.Key]; ok {
			return fmt.Errorf("Duplicate key names for: %v", item.Key)
		}
		itemMap[item.Key] = item
	}

	sm.ItemMap = itemMap
	return nil
}

func (sm *SheetsMap) GetItem(key string) (Item, error) {
	if item, ok := sm.ItemMap[key]; !ok {
		return Item{}, fmt.Errorf("Cannot find key %v", key)
	} else {
		return item, nil
	}
}

func (sm *SheetsMap) GetOrCreateItem(itemKey string) (Item, error) {
	itemKey = TrimSpaceToLower(itemKey)
	if item, ok := sm.ItemMap[itemKey]; !ok {
		item := Item{
			Key:  itemKey,
			Data: map[string]string{},
		}
		if len(sm.Columns) > 0 {
			item.Data[sm.Columns[0].Value] = itemKey
		}

		itemCount := len(sm.ItemMap)
		nextRowIdx := itemCount + 1
		item.Row = uint(nextRowIdx)

		sm.Sheet.Update(nextRowIdx, 0, itemKey)
		err := sm.Sheet.Synchronize()
		if err == nil {
			sm.ItemMap[itemKey] = item
		}
		return item, err
	} else {
		return item, nil
	}
}

func (sm *SheetsMap) UpdateItem(item Item, key, val string, synchronize bool) (string, error) {
	// Get key column
	keyLc := TrimSpaceToLower(key)
	col, ok := sm.ColumnMapKeyLc[keyLc]
	if !ok {
		return "", fmt.Errorf("Key Not Found: %v", key)
	}

	// Process value
	str, err := col.ValueToCanonical(val)
	if err != nil {
		return str, err
	}

	item.Data[keyLc] = str
	return "", sm.SynchronizeItem(item)
}

func (sm *SheetsMap) SynchronizeItem(item Item) error {
	rowIdx := item.Row
	for colIdx, col := range sm.Columns {
		if val, ok := item.Data[col.Value]; ok {
			sm.Sheet.Update(int(rowIdx), colIdx, val)
		} else {
			sm.Sheet.Update(int(rowIdx), colIdx, "")
		}
	}
	return sm.Sheet.Synchronize()
}

func (sm *SheetsMap) SetItemKeyColValue(itemKey, colKeyRaw, colValRaw string) (Item, error) {
	item, err := sm.GetOrCreateItem(itemKey)
	if err != nil {
		return item, err
	}

	colKey := strings.TrimSpace(colKeyRaw)
	colKeyLc := strings.ToLower(colKeyRaw)
	col, ok := sm.ColumnMapKeyLc[colKeyLc]
	if !ok {
		return item, fmt.Errorf("Column not found [%v]", colKey)
	}

	colVal, err := col.ValueToCanonical(colValRaw)
	if err != nil {
		return item, err
	}

	item.Data[colKey] = colVal
	sm.ItemMap[itemKey] = item
	return item, nil
}

type Intent struct {
	Name  string
	Slots map[string]string
}

func TrimSpaceToLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (sm *SheetsMap) SetItemKeyString(itemKey, cmdRaw string) (Intent, error) {
	cmdRawLc := TrimSpaceToLower(cmdRaw)
	intent := Intent{Slots: map[string]string{}}

	fmt.Println(len(sm.Columns))

	for _, col := range sm.ColumnMapKeyLc {
		colKeyLc := TrimSpaceToLower(col.Value)
		pat := fmt.Sprintf("^%v\\s*(.*)$", colKeyLc)
		pat1, err := regexp.Compile(pat)
		if err != nil {
			return intent, fmt.Errorf("Cannot compile regexp for colKey")
		}
		m := pat1.FindStringSubmatch(cmdRawLc)
		if len(m) == 2 {
			valCanonical, err := col.ValueToCanonical(m[1])
			if err != nil {
				return intent, fmt.Errorf("E_INCORRECT_VALUE")
			}
			item, err := sm.SetItemKeyColValue(itemKey, col.Value, valCanonical)
			if err != nil {
				panic(err)
			}
			err = sm.SynchronizeItem(item)
			if err != nil {
				panic(err)
			}
			break
		}
	}
	return intent, nil

}
