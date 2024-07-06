package transform

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"pkg/common/errors"
)

// Transformer is responsible for data transformations
type Transformer struct{}

// NewTransformer creates a new Transformer
// ```
// Example:
// 		transformer := transform.NewTransformer()

//     // Example map data
//     data := map[string]interface{}{
//         "name": "John Doe",
//         "age":  30,
//         "address": map[string]interface{}{
//             "street": "123 Main St",
//             "city":   "New York",
//             "zip":    "10001",
//         },
//         "hobbies": []string{"reading", "swimming"},
//     }

//     // Convert to struct
//     result, err := transformer.TransformData(data, "struct")
//     if err != nil {
//         log.Fatalf("Failed to convert to struct: %v", err)
//     }

//     // Print the resulting struct
//     fmt.Printf("%+v\n", result)

//     // You can now access fields like this:
//     resultValue := reflect.ValueOf(result)
//     fmt.Printf("Name: %v\n", resultValue.FieldByName("Name").Interface())
//     fmt.Printf("Age: %v\n", resultValue.FieldByName("Age").Interface())

//	 // Access nested struct
//	 address := resultValue.FieldByName("Address").Interface()
//	 addressValue := reflect.ValueOf(address)
//	 fmt.Printf("City: %v\n", addressValue.FieldByName("City").Interface())
//	fmt.Printf("Zip: %v\n", addressValue.FieldByName("Zip").Interface())
//
// ```
func NewTransformer() *Transformer {
	return &Transformer{}
}

// TransformData converts data from one format to another
func (t *Transformer) TransformData(data interface{}, targetFormat string) (interface{}, error) {
	switch targetFormat {
	case "json":
		return t.toJSON(data)
	case "map":
		return t.toMap(data)
	case "struct":
		return t.toStruct(data)
	case "array":
		return t.toArray(data)
	default:
		return nil, errors.NewError(errors.ErrorTypeUnsupported, fmt.Sprintf("unsupported target format: %s", targetFormat), nil)
	}
}

// toStruct converts data to a struct using reflection
func (t *Transformer) toStruct(data interface{}) (interface{}, error) {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	switch dataValue.Kind() {
	case reflect.Struct:
		return data, nil // Already a struct
	case reflect.Map:
		return t.mapToStruct(dataValue)
	default:
		return nil, errors.NewError(errors.ErrorTypeTransformation, "unsupported data type for struct conversion", nil)
	}
}

// mapToStruct converts a map to a struct
func (t *Transformer) mapToStruct(mapValue reflect.Value) (interface{}, error) {
	if mapValue.Kind() != reflect.Map {
		return nil, errors.NewError(errors.ErrorTypeTransformation, "input is not a map", nil)
	}

	structType := reflect.StructOf(t.mapToStructFields(mapValue))
	structValue := reflect.New(structType).Elem()

	for _, key := range mapValue.MapKeys() {
		fieldName := t.normalizeFieldName(key.String())
		fieldValue := mapValue.MapIndex(key)

		if structValue.FieldByName(fieldName).IsValid() {
			if err := t.setStructField(structValue.FieldByName(fieldName), fieldValue); err != nil {
				return nil, err
			}
		}
	}

	return structValue.Interface(), nil
}

// mapToStructFields creates struct fields from a map
func (t *Transformer) mapToStructFields(mapValue reflect.Value) []reflect.StructField {
	var fields []reflect.StructField

	for _, key := range mapValue.MapKeys() {
		fieldName := t.normalizeFieldName(key.String())
		fieldValue := mapValue.MapIndex(key)
		fieldType := fieldValue.Type()

		// Handle nested maps
		if fieldValue.Kind() == reflect.Map {
			nestedFields := t.mapToStructFields(fieldValue)
			fieldType = reflect.StructOf(nestedFields)
		}

		fields = append(fields, reflect.StructField{
			Name: fieldName,
			Type: fieldType,
		})
	}

	return fields
}

// setStructField sets the value of a struct field
func (t *Transformer) setStructField(field reflect.Value, value reflect.Value) error {
	if !field.CanSet() {
		return errors.NewError(errors.ErrorTypeTransformation, "cannot set field value", nil)
	}

	if field.Kind() == reflect.Struct && value.Kind() == reflect.Map {
		nestedStruct, err := t.mapToStruct(value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(nestedStruct))
		return nil
	}

	if field.Type() != value.Type() {
		if value.Type().ConvertibleTo(field.Type()) {
			value = value.Convert(field.Type())
		} else {
			return errors.NewError(errors.ErrorTypeTransformation, "incompatible types for field assignment", nil)
		}
	}

	field.Set(value)
	return nil
}

// normalizeFieldName converts a string to a valid Go struct field name
func (t *Transformer) normalizeFieldName(name string) string {
	// Capitalize the first letter
	name = strings.Title(name)
	// Remove any characters that are not letters, numbers, or underscores
	name = strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, name)
	// Ensure the field name starts with a letter
	if len(name) > 0 && (name[0] >= '0' && name[0] <= '9') {
		name = "F" + name
	}
	return name
}

// toJSON converts data to JSON format
func (t *Transformer) toJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", errors.NewError(errors.ErrorTypeTransformation, "failed to convert data to JSON", err)
	}
	return string(jsonData), nil
}

// toMap converts data to a map[string]interface{}
func (t *Transformer) toMap(data interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		return v, nil
	case string:
		if err := json.Unmarshal([]byte(v), &result); err != nil {
			return nil, errors.NewError(errors.ErrorTypeTransformation, "failed to convert string to map", err)
		}
		return result, nil
	default:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, errors.NewError(errors.ErrorTypeTransformation, "failed to marshal data to JSON", err)
		}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			return nil, errors.NewError(errors.ErrorTypeTransformation, "failed to unmarshal JSON to map", err)
		}
		return result, nil
	}
}

// toArray converts data to an array ([]interface{})
func (t *Transformer) toArray(data interface{}) ([]interface{}, error) {
	switch v := data.(type) {
	case []interface{}:
		return v, nil
	case string:
		var result []interface{}
		if err := json.Unmarshal([]byte(v), &result); err != nil {
			return nil, errors.NewError(errors.ErrorTypeTransformation, "failed to convert string to array", err)
		}
		return result, nil
	default:
		value := reflect.ValueOf(data)
		if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
			return nil, errors.NewError(errors.ErrorTypeTransformation, "data is not a slice or array", nil)
		}
		result := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			result[i] = value.Index(i).Interface()
		}
		return result, nil
	}
}

// ConvertType converts a value to the specified type
func (t *Transformer) ConvertType(value interface{}, targetType string) (interface{}, error) {
	switch targetType {
	case "string":
		return fmt.Sprintf("%v", value), nil
	case "int":
		switch v := value.(type) {
		case string:
			return strconv.Atoi(v)
		case float64:
			return int(v), nil
		default:
			return 0, errors.NewError(errors.ErrorTypeTransformation, "unable to convert to int", nil)
		}
	case "float":
		switch v := value.(type) {
		case string:
			return strconv.ParseFloat(v, 64)
		case int:
			return float64(v), nil
		default:
			return 0.0, errors.NewError(errors.ErrorTypeTransformation, "unable to convert to float", nil)
		}
	case "bool":
		switch v := value.(type) {
		case string:
			return strconv.ParseBool(v)
		case int:
			return v != 0, nil
		default:
			return false, errors.NewError(errors.ErrorTypeTransformation, "unable to convert to bool", nil)
		}
	case "time":
		switch v := value.(type) {
		case string:
			return time.Parse(time.RFC3339, v)
		case int64:
			return time.Unix(v, 0), nil
		default:
			return time.Time{}, errors.NewError(errors.ErrorTypeTransformation, "unable to convert to time", nil)
		}
	default:
		return nil, errors.NewError(errors.ErrorTypeUnsupported, fmt.Sprintf("unsupported target type: %s", targetType), nil)
	}
}

// FlattenMap flattens a nested map structure
func (t *Transformer) FlattenMap(data map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch child := v.(type) {
		case map[string]interface{}:
			childMap := t.FlattenMap(child, key)
			for ck, cv := range childMap {
				result[ck] = cv
			}
		case []interface{}:
			for i, item := range child {
				if childMap, ok := item.(map[string]interface{}); ok {
					childFlat := t.FlattenMap(childMap, fmt.Sprintf("%s[%d]", key, i))
					for ck, cv := range childFlat {
						result[ck] = cv
					}
				} else {
					result[fmt.Sprintf("%s[%d]", key, i)] = item
				}
			}
		default:
			result[key] = v
		}
	}
	return result
}

// UnflattenMap reverses the flattening process
func (t *Transformer) UnflattenMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		parts := strings.Split(k, ".")
		m := result
		for _, part := range parts[:len(parts)-1] {
			if _, ok := m[part]; !ok {
				m[part] = make(map[string]interface{})
			}
			m = m[part].(map[string]interface{})
		}
		m[parts[len(parts)-1]] = v
	}
	return result
}
