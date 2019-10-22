package table

import (
	"io/ioutil"
	"reflect"
	"strings"

	"utils"
)

// Table to hold CSV-formatted table
type Table struct {
	colNames []string
	colTypes map[string]string

	bodyRows [][]interface{}
}

// Load to load table either from a file which defines the
// table or a string which defines the table.
// The parameter tableStr could be either a path to the
// file or a string which defines the table.
func (tbl *Table) Load(tableStr string) error {
	buf, err := ioutil.ReadFile(tableStr)
	if err == nil {
		return tbl.LoadFromCsvString(string(buf))
	}

	return tbl.LoadFromCsvString(tableStr)
}

// LoadFromCsvString to load CSV-formatted table string.
// The table should follow the pattern:
//   `
//   nameCol1,nameCol2,nameCol3
//   typeCol1,typeCol2,typeCol3
//   valCol11,valCol21,valCol31
//   `
//
// For example:
//   `
//   client,region,domestic,total,compliance
//   string,string,bool,int,float
//   TT,APAC,false,100,0.9
//   ...
//   `
func (tbl *Table) LoadFromCsvString(tableStr string) error {
	// TODO: validate the input and error handling

	tableStr = strings.TrimSpace(tableStr)
	rows := strings.Split(tableStr, "\n")

	// Load column names from the first row
	tbl.colNames = utils.SplitAndTrimSpace(rows[0], ",")

	// Load column types from the second row
	tbl.loadColumnTypes(rows[1])

	// Load all data (body) rows
	tbl.loadBodyRows(rows[2:]...)

	return nil
}

func (tbl *Table) GetColumnNameByIndex(i int) string {
	// TODO: validate input
	return tbl.colNames[i]
}

func (tbl *Table) GetColumnTypeByName(colName string) string {
	// TODO: validate input
	return tbl.colTypes[colName]
}

// GetRowByColumns to query the row which contains the same values of the fields
// passed in (colNameToValue).
// For example, the passed-in map with the value "client:TT, region:APAC" will
// return the row "TT,APAC,false,100,0.9" of the table given below:
//   `
//   client,region,domestic,total,compliance
//   string,string,bool,int,float
//   TT,APAC,false,100,0.9
//   ATPIUK,EU,true,99,0.8
//   TT,APAC,true,66,0.95
//   ...
//   `
// Note if there are multiple rows matched, the first occurrance is returned. In our
// above example, the "TT,APAC,false,100,0.9" rather than "TT,APAC,true,66,0.95" is returned.
func (tbl *Table) GetRowByColumns(colNameToValue map[string]interface{}) map[string]interface{} {
	//TODO exception handling

	for _, row := range tbl.bodyRows {
		cntMatched := 0
		for i := range row {
			colName := tbl.colNames[i]
			val, ok := colNameToValue[colName]
			if ok {
				if !reflect.DeepEqual(val, row[i]) {
					break
				}

				cntMatched++
			}
		}

		if cntMatched == len(colNameToValue) {
			rowMatched := make(map[string]interface{})
			for i, val := range row {
				rowMatched[tbl.colNames[i]] = val
			}

			return rowMatched
		}
	}

	return nil
}

func (tbl *Table) loadColumnTypes(typesStr string) {
	colTypes := utils.SplitAndTrimSpace(typesStr, ",")

	tbl.colTypes = make(map[string]string)
	for i, col := range tbl.colNames {
		tbl.colTypes[col] = colTypes[i]
	}
}

func (tbl *Table) loadBodyRows(bodyRows ...string) {
	tbl.bodyRows = make([][]interface{}, len(bodyRows))

	for i, row := range bodyRows {
		tbl.bodyRows[i] = tbl.loadSingleRow(row)
	}
}

func (tbl *Table) loadSingleRow(row string) []interface{} {
	fieldsInStr := utils.SplitAndTrimSpace(row, ",")
	fields := make([]interface{}, len(fieldsInStr))

	for i, valStr := range fieldsInStr {
		valType := tbl.colTypes[tbl.colNames[i]]
		fields[i], _ = utils.ConvertStringToType(valStr, valType)
	}

	return fields
}
