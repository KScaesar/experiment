package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func main() {
	data := []byte(`
{
  "version": "1.0",
  "rules": [
    {
      "resource": {
        "path": "/api/data/documents"
      },
      "allowOrigins": [
        "http://this.example.com",
        "http://that.example.com"
      ],
      "allowMethods": [
        "GET"
      ],
      "allowCredentials": true
    }
  ],
  "newFinalTest": null
}`)
	payload := make(map[string]interface{})
	json.Unmarshal(data, &payload)

	newPayload := jsonConvertObject(payload)
	newData, _ := json.MarshalIndent(newPayload, "", " ")
	fmt.Println(string(newData))
}

func jsonConvertSlice(payload []interface{}) []interface{} {
	newPayload := make([]interface{}, len(payload))

	for index, value := range payload {
		var newValue interface{}

		switch reflect.ValueOf(value).Kind() {
		case reflect.Slice:
			newValue = jsonConvertSlice(value.([]interface{}))

		case reflect.Map:
			newValue = jsonConvertObject(value.(map[string]interface{}))

		default:
			newValue = value
		}

		newPayload[index] = newValue
	}

	return newPayload
}

func jsonConvertObject(payload map[string]interface{}) map[string]interface{} {
	newPayload := make(map[string]interface{}, len(payload))

	for key, value := range payload {
		newKey := toSnake(key)
		var newValue interface{}

		switch reflect.ValueOf(value).Kind() {
		case reflect.Slice:
			newValue = jsonConvertSlice(value.([]interface{}))

		case reflect.Map:
			newValue = jsonConvertObject(value.(map[string]interface{}))

		default:
			newValue = value
		}

		newPayload[newKey] = newValue
	}

	return newPayload
}

func toSnake(text string) string {
	result := strings.Builder{}
	lower, upper := 'a', 'A'
	diff := lower - upper

	for _, char := range text {
		if char >= 'A' && char <= 'Z' {
			result.WriteRune('_')
			result.WriteRune(char + diff)
			continue
		}
		result.WriteRune(char)
	}

	return result.String()
}
