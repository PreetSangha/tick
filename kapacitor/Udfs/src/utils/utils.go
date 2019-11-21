package utils

import (
	"strconv"
	"strings"
)

// ConvertStringToType to convert string to the type specified and then
// return the converted value back as interface{}.
// If the type is not supported, nil will be returned back.
func ConvertStringToType(value string, valType string) (interface{}, error) {
	var result interface{}
	var err error

	switch valType {
	case "string":
		result = value
	case "int":
		tempVal, tempErr := strconv.Atoi(value)
		err = tempErr
		result = int64(tempVal)
	case "float":
		result, err = strconv.ParseFloat(value, 64)
	case "bool":
		result, err = strconv.ParseBool(value)
	}

	return result, err
}

// SplitAndTrimSpace to split a string by the delemiter passed-in
// and then trim space.
func SplitAndTrimSpace(rawStr string, splitBy string) []string {
	rawStr = strings.TrimSpace(rawStr)
	if len(rawStr) == 0 {
		return []string{}
	}

	items := strings.Split(rawStr, splitBy)

	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	return items
}
