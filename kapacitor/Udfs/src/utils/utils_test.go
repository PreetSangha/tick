package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConvertStringToType(t *testing.T) {
	// test cases
	for _, tc := range [...]struct {
		value    string
		valType  string
		expected interface{}
	}{
		{"TT", "string", "TT"},
		{"100", "int", int64(100)},
		{"999.999", "float", 999.999},
		{"101.00", "float", 101.0},
		{"true", "bool", true},
		{"false", "bool", false},
		{"Unknown", "unknown", nil},
	} {
		t.Run(fmt.Sprintf("Convert string '%s' to type '%s'", tc.value, tc.valType), func(t *testing.T) {
			actual, _ := ConvertStringToType(tc.value, tc.valType)
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func TestSplitAndTrimSpace(t *testing.T) {
	type checkFunc func([]string) error

	// checkers
	isEmptyList := func(have []string) error {
		if len(have) > 0 {
			return fmt.Errorf("expected empty list, found %v", have)
		}
		return nil
	}
	is := func(want ...string) checkFunc {
		return func(have []string) error {
			if !reflect.DeepEqual(have, want) {
				return fmt.Errorf("expected list %v, found %v", want, have)
			}
			return nil
		}
	}

	// test cases
	for _, tc := range [...]struct {
		strToSplit string
		splitBy    string
		check      checkFunc
	}{
		{"", ",", isEmptyList},
		{"client,domain,slo", ",", is("client", "domain", "slo")},
		{"client,domain,slo", ":", is("client,domain,slo")},
		{" T T : 123: 11.99:11.0 :false ", ":", is("T T", "123", "11.99", "11.0", "false")},
		{"\tT T , 123,11.99,11.0 ,false\r\n,", ",", is("T T", "123", "11.99", "11.0", "false", "")},
		{`client,
		domain,
		slo`, ",", is("client", "domain", "slo")},
	} {
		t.Run(fmt.Sprintf("Split %s by '%s' and trim space", tc.strToSplit, tc.splitBy), func(t *testing.T) {
			result := SplitAndTrimSpace(tc.strToSplit, tc.splitBy)
			if err := tc.check(result); err != nil {
				t.Error(err)
			}
		})
	}
}
