package table

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func getTable() *Table {
	tableCsv := `client,region,domestic,total,compliance
	string,string,bool,int,float
	TT,APAC,false,100,0.9
	ATPIUK,EU,true,99,0.8
	TT,APAC,true,66,0.95	`

	t := Table{}
	t.LoadFromCsvString(tableCsv)

	return &t
}

func TestLoad(t *testing.T) {
	tableFilePath, _ := filepath.Abs("./testdata/table.csv")
	// test cases
	for _, tc := range [...]struct {
		in       string
		table    *Table
		expected *Table
	}{
		{`client,region,domestic,total,compliance
		string,string,bool,int,float
		TT,APAC,false,100,0.9
		ATPIUK,EU,true,99,0.8
		TT,APAC,true,66,0.95`,
			&Table{},
			&Table{
				[]string{"client", "region", "domestic", "total", "compliance"},
				map[string]string{"client": "string", "region": "string", "total": "int", "compliance": "float", "domestic": "bool"},
				[][]interface{}{{"TT", "APAC", false, int64(100), 0.9}, {"ATPIUK", "EU", true, int64(99), 0.8}, {"TT", "APAC", true, int64(66), 0.95}},
			},
		},
		{tableFilePath,
			&Table{},
			&Table{
				[]string{"client", "region", "domestic", "total", "compliance"},
				map[string]string{"client": "string", "region": "string", "total": "int", "compliance": "float", "domestic": "bool"},
				[][]interface{}{{"TT", "APAC", false, int64(100), 0.9}, {"ATPIUK", "EU", true, int64(99), 0.8}, {"TT", "APAC", true, int64(66), 0.95}},
			},
		},
	} {
		t.Run(fmt.Sprintf("Load table from CSV-formatted string or file %s", tc.in), func(t *testing.T) {
			tc.table.Load(tc.in)
			if !reflect.DeepEqual(tc.expected.colNames, tc.table.colNames) {
				t.Errorf("expected colNames %v, actual colNames %v", tc.expected.colNames, tc.table.colNames)
			}
			if !reflect.DeepEqual(tc.expected.colTypes, tc.table.colTypes) {
				t.Errorf("expected colTypes %v, actual colTypes %v", tc.expected.colTypes, tc.table.colTypes)
			}
			if !reflect.DeepEqual(tc.expected.bodyRows, tc.table.bodyRows) {
				t.Errorf("expected bodyRows %v, actual bodyRows %v", tc.expected.bodyRows, tc.table.bodyRows)
			}
		})
	}
}

func TestLoadFromCsvString(t *testing.T) {
	// TODO: add negative test cases

	// test cases
	for _, tc := range [...]struct {
		tableCsvStr string
		table       *Table
		expected    *Table
	}{
		{"client,total,compliance,domestic\nstring,int,float,bool\nTT,100,0.9,true\nFCL,200,0.95,false",
			&Table{},
			&Table{
				[]string{"client", "total", "compliance", "domestic"},
				map[string]string{"client": "string", "total": "int", "compliance": "float", "domestic": "bool"},
				[][]interface{}{{"TT", int64(100), 0.9, true}, {"FCL", int64(200), 0.95, false}},
			},
		},
		{`client,total, compliance,domestic
		string,int,float, bool
		TT,100, 0.9,true
		FCL,200,0.95, false`,
			&Table{},
			&Table{
				[]string{"client", "total", "compliance", "domestic"},
				map[string]string{"client": "string", "total": "int", "compliance": "float", "domestic": "bool"},
				[][]interface{}{{"TT", int64(100), 0.9, true}, {"FCL", int64(200), 0.95, false}},
			},
		},
	} {
		t.Run(fmt.Sprintf("Load table from CSV-formatted string %s", tc.tableCsvStr), func(t *testing.T) {
			tc.table.LoadFromCsvString(tc.tableCsvStr)
			if !reflect.DeepEqual(tc.expected.colNames, tc.table.colNames) {
				t.Errorf("expected colNames %v, actual colNames %v", tc.expected.colNames, tc.table.colNames)
			}
			if !reflect.DeepEqual(tc.expected.colTypes, tc.table.colTypes) {
				t.Errorf("expected colTypes %v, actual colTypes %v", tc.expected.colTypes, tc.table.colTypes)
			}
			if !reflect.DeepEqual(tc.expected.bodyRows, tc.table.bodyRows) {
				t.Errorf("expected bodyRows %v, actual bodyRows %v", tc.expected.bodyRows, tc.table.bodyRows)
			}
		})
	}
}

func TestGetRowByColumns(t *testing.T) {
	// TODO: add negative test cases
	tbl := getTable()

	// test cases
	for _, tc := range [...]struct {
		query    map[string]interface{}
		table    *Table
		expected map[string]interface{}
	}{
		{map[string]interface{}{"client": "TT", "region": "APAC"},
			tbl,
			map[string]interface{}{"client": "TT", "region": "APAC", "domestic": false, "total": int64(100), "compliance": 0.9},
		},
		{map[string]interface{}{"client": "Unknown", "region": "APAC"},
			tbl,
			nil,
		},
	} {
		t.Run(fmt.Sprintf("Get table row by cloumns %v", tc.query), func(t *testing.T) {
			actual := tc.table.GetRowByColumns(tc.query)
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("expected colNames %v, actual colNames %v", tc.expected, actual)
			}
		})
	}
}
