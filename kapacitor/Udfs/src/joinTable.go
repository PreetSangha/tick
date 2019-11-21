package main

import (
	"log"
	"os"

	"table"
	"utils"

	"github.com/influxdata/kapacitor/udf/agent"
)

type joinTableHandler struct {
	table        *table.Table
	on           []string
	defaultValue map[string]interface{}

	agent *agent.Agent
}

func newJoinTableHandler(agent *agent.Agent) *joinTableHandler {
	return &joinTableHandler{
		agent: agent,
		table: &table.Table{},
	}
}

// Return the InfoResponse. Describing the properties of this UDF agent.
func (*joinTableHandler) Info() (*agent.InfoResponse, error) {
	info := &agent.InfoResponse{
		Wants:    agent.EdgeType_BATCH,
		Provides: agent.EdgeType_BATCH,

		Options: map[string]*agent.OptionInfo{
			"table":   {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
			"on":      {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
			"default": {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
		},
	}

	return info, nil
}

// Initialze the handler based of the provided options.
func (jt *joinTableHandler) Init(r *agent.InitRequest) (*agent.InitResponse, error) {
	init := &agent.InitResponse{
		Success: true,
		Error:   "",
	}

	for _, opt := range r.Options {
		switch opt.Name {
		case "table":
			jt.table.Load(opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue)
		case "on":
			jt.on = utils.SplitAndTrimSpace(opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue, ",")
		case "default":
			jt.parseDefault(opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue)
		}
	}

	if jt.table == nil || len(jt.on) == 0 {
		init.Success = false
		init.Error = "must supply 'table' and 'on'"
	}

	return init, nil
}

func (jt *joinTableHandler) parseDefault(defaultStr string) {
	defaulVals := parseKeyValuePairString(defaultStr)

	res := make(map[string]interface{})
	for k, v := range defaulVals {
		res[k], _ = utils.ConvertStringToType(v, jt.table.GetColumnTypeByName(k))
	}

	jt.defaultValue = res
}

func parseKeyValuePairString(kvPairStr string) map[string]string {
	// Parse string like "client: TT, domain: zeno, totalVolue: 100"
	// to a map structure.

	res := make(map[string]string)

	vals := utils.SplitAndTrimSpace(kvPairStr, ",")
	for _, val := range vals {
		fields := utils.SplitAndTrimSpace(val, ":")
		res[fields[0]] = fields[1]
	}

	return res
}

// Create a snapshot of the running state of the process.
func (*joinTableHandler) Snapshot() (*agent.SnapshotResponse, error) {
	return &agent.SnapshotResponse{}, nil
}

// Restore a previous snapshot.
func (*joinTableHandler) Restore(req *agent.RestoreRequest) (*agent.RestoreResponse, error) {
	return &agent.RestoreResponse{
		Success: true,
	}, nil
}

// Start working with the next batch
func (*joinTableHandler) BeginBatch(begin *agent.BeginBatch) error {
	return nil
}

func (jt *joinTableHandler) getPointFieldByName(name string, p *agent.Point) interface{} {
	//TODO exception handling

	var val interface{}
	ok := false

	switch jt.table.GetColumnTypeByName(name) {
	case "string":
		if val, ok = p.Tags[name]; !ok {
			val, ok = p.FieldsString[name]
		}
	case "int":
		val, ok = p.FieldsInt[name]
	case "float":
		val, ok = p.FieldsDouble[name]
	case "bool":
		val, ok = p.FieldsBool[name]
	}

	return val
}

func (jt *joinTableHandler) Point(p *agent.Point) error {
	//TODO exception handling

	// Create a map to hold the name and its value of columns (fields)
	// that we are going to join on.
	piontFields := make(map[string]interface{})
	for _, colName := range jt.on {
		piontFields[colName] = jt.getPointFieldByName(colName, p)
	}

	// Get the row from the table that matches the current data point.
	// If no row is found and a default value is provided, then it joins the
	// data using the default value.
	row := jt.table.GetRowByColumns(piontFields)
	if row == nil {
		row = jt.defaultValue
	}

	joinPointWithRowOnFields(p, row, piontFields)

	// Send the new data point (after joined) back to Kapacitor
	jt.agent.Responses <- &agent.Response{
		Message: &agent.Response_Point{
			Point: p,
		},
	}

	return nil
}

func joinPointWithRowOnFields(p *agent.Point, row map[string]interface{}, onFields map[string]interface{}) {
	for colName, val := range row {
		if _, existing := onFields[colName]; existing {
			continue
		}

		switch val.(type) {
		case string:
			if p.FieldsString == nil {
				p.FieldsString = make(map[string]string)
			}
			p.FieldsString[colName] = val.(string)
		case int64:
			if p.FieldsInt == nil {
				p.FieldsInt = make(map[string]int64)
			}
			p.FieldsInt[colName] = val.(int64)
		case float64:
			if p.FieldsDouble == nil {
				p.FieldsDouble = make(map[string]float64)
			}
			p.FieldsDouble[colName] = val.(float64)
		case bool:
			if p.FieldsBool == nil {
				p.FieldsBool = make(map[string]bool)
			}
			p.FieldsBool[colName] = val.(bool)
		}
	}
}

func (jt *joinTableHandler) EndBatch(end *agent.EndBatch) error {
	return nil
}

// Stop the handler gracefully.
func (jt *joinTableHandler) Stop() {
	close(jt.agent.Responses)
}

func main() {
	a := agent.New(os.Stdin, os.Stdout)
	h := newJoinTableHandler(a)
	a.Handler = h

	log.Println("Starting agent")
	a.Start()
	err := a.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
