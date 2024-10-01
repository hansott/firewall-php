package api_discovery

import (
	. "main/aikido_types"
	"main/utils"
	"reflect"
)

// mergeDataSchemas merges two DataSchema objects.
func mergeDataSchemas(first *DataSchema, second *DataSchema) *DataSchema {
	result := first

	// Cannot merge different types
	if !isSameType(first.Type, second.Type) {
		return mergeTypes(first, second)
	}

	if first.Properties != nil && second.Properties != nil {
		result.Properties = make(map[string]*DataSchema)
		for key, value := range first.Properties {
			result.Properties[key] = value
		}

		for key, secondProp := range second.Properties {
			if firstProp, ok := result.Properties[key]; ok {
				result.Properties[key] = mergeDataSchemas(firstProp, secondProp)
			} else {
				result.Properties[key] = secondProp
				// Mark as optional since it's only present in the second schema
				result.Properties[key].Optional = utils.BoolPointer(true)
			}
		}

		for key := range first.Properties {
			if _, ok := second.Properties[key]; !ok {
				result.Properties[key].Optional = utils.BoolPointer(true)
			}
		}
	}

	if first.Items != nil && second.Items != nil {
		result.Items = mergeDataSchemas(first.Items, second.Items)
	}

	return result
}

// isSameType checks if both types are the same or compatible.
func isSameType(first, second interface{}) bool {
	// Check for arrays of types
	if reflect.TypeOf(first) != reflect.TypeOf(second) {
		return false
	}

	switch firstVal := first.(type) {
	case string:
		return firstVal == second.(string)
	case []string:
		return doTypeArraysMatch(firstVal, second.([]string))
	}
	return false
}

// doTypeArraysMatch compares two arrays of types, ignoring the order.
func doTypeArraysMatch(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}

	for _, typ := range first {
		found := false
		for _, typ2 := range second {
			if typ == typ2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// mergeTypes merges two DataSchema objects of different types.
func mergeTypes(first *DataSchema, second *DataSchema) *DataSchema {
	// Cannot merge arrays and objects or primitives with non-primitives
	if !onlyContainsPrimitiveTypes(first.Type) || !onlyContainsPrimitiveTypes(second.Type) {
		// Prefer non-null type
		if first.Type == "null" {
			return second
		}
		return first
	}

	first.Type = mergeTypeArrays(first.Type, second.Type)
	return first
}

// mergeTypeArrays merges two types into a single array of unique types.
func mergeTypeArrays(first, second interface{}) interface{} {
	var firstArr, secondArr []string

	switch v := first.(type) {
	case string:
		firstArr = []string{v}
	case []string:
		firstArr = v
	}

	switch v := second.(type) {
	case string:
		secondArr = []string{v}
	case []string:
		secondArr = v
	}

	typeSet := make(map[string]bool)
	for _, t := range append(firstArr, secondArr...) {
		typeSet[t] = true
	}

	merged := make([]string, 0, len(typeSet))
	for t := range typeSet {
		merged = append(merged, t)
	}

	return merged
}
