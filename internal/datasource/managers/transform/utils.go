package transform

import (
	"fmt"
	"strings"

	"pkg/common/errors"
)

// ExtractField extracts a nested field from a map using dot notation
func ExtractField(data map[string]interface{}, field string) (interface{}, error) {
	parts := strings.Split(field, ".")
	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			if val, ok := current[part]; ok {
				return val, nil
			}
			return nil, errors.NewError(errors.ErrorTypeNotFound, "field not found", nil)
		}
		if val, ok := current[part]; ok {
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil, errors.NewError(errors.ErrorTypeTransformation, "invalid path", nil)
			}
		} else {
			return nil, errors.NewError(errors.ErrorTypeNotFound, "field not found", nil)
		}
	}
	return nil, errors.NewError(errors.ErrorTypeTransformation, "invalid path", nil)
}

// SetField sets a nested field in a map using dot notation
func SetField(data map[string]interface{}, field string, value interface{}) error {
	parts := strings.Split(field, ".")
	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			current[part] = value
			return nil
		}
		if val, ok := current[part]; ok {
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				newMap := make(map[string]interface{})
				current[part] = newMap
				current = newMap
			}
		} else {
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
		}
	}
	return errors.NewError(errors.ErrorTypeTransformation, "invalid path", nil)
}

// MergeMap merges two maps, with the second map taking precedence
func MergeMap(m1, m2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		if v1, ok := result[k]; ok {
			if map1, isMap1 := v1.(map[string]interface{}); isMap1 {
				if map2, isMap2 := v.(map[string]interface{}); isMap2 {
					result[k] = MergeMap(map1, map2)
					continue
				}
			}
		}
		result[k] = v
	}
	return result
}

// GroupBy groups a slice of maps by a specified key
func GroupBy(data []map[string]interface{}, key string) (map[string][]map[string]interface{}, error) {
	result := make(map[string][]map[string]interface{})
	for _, item := range data {
		value, err := ExtractField(item, key)
		if err != nil {
			return nil, err
		}
		strValue := fmt.Sprintf("%v", value)
		result[strValue] = append(result[strValue], item)
	}
	return result, nil
}

// FilterSlice filters a slice of maps based on a condition function
func FilterSlice(data []map[string]interface{}, condition func(map[string]interface{}) bool) []map[string]interface{} {
	var result []map[string]interface{}
	for _, item := range data {
		if condition(item) {
			result = append(result, item)
		}
	}
	return result
}
