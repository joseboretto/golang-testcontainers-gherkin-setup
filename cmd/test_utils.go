package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// compareJSON compares two JSON strings and returns true if they are equal.
func compareJSON(jsonStr1, jsonStr2 string) (bool, error) {
	var obj1, obj2 map[string]interface{}

	// Unmarshal the first JSON string
	if err := json.Unmarshal([]byte(jsonStr1), &obj1); err != nil {
		return false, fmt.Errorf("error unmarshalling jsonStr1: %v", err)
	}

	// Unmarshal the second JSON string
	if err := json.Unmarshal([]byte(jsonStr2), &obj2); err != nil {
		return false, fmt.Errorf("error unmarshalling jsonStr2: %v", err)
	}
	// Compare the two maps
	return reflect.DeepEqual(obj1, obj2), nil
}

func compareJSONArrays(jsonStr1, jsonStr2 string) (bool, error) {
	var obj1, obj2 []map[string]interface{}

	// Unmarshal the first JSON string
	if err := json.Unmarshal([]byte(jsonStr1), &obj1); err != nil {
		return false, fmt.Errorf("error unmarshalling jsonStr1: %v", err)
	}

	// Unmarshal the second JSON string
	if err := json.Unmarshal([]byte(jsonStr2), &obj2); err != nil {
		return false, fmt.Errorf("error unmarshalling jsonStr2: %v", err)
	}
	// Compare the two maps
	return reflect.DeepEqual(obj1, obj2), nil
}
