package tagconv

import (
	"fmt"
	"github.com/imdario/mergo"
	"reflect"
	"strings"
)

var tagName = ""

// getMapOfAllKeyValues builds a map of the fully specified key and the value from the struct tag
// eg:
/*
	"data.call": 2,
	"data.text": "2",
	"email": "3",
	"hello": "1",
	"id": 1,
	"name": "2",
	"object.data.world": "6",
	"object.name": "4",
	"object.text": "5"
*/
func getMapOfAllKeyValues(s interface{}) (*map[string]interface{}, error) {
	var vars = make(map[string]interface{}) // this will hold the variables as a map (JSON)

	// TODO: catch panics when reflecting unexported fields
	// TODO: handle arrays

	// get value of object
	t := reflect.ValueOf(s)
	if t.IsZero() {
		return nil, fmt.Errorf("empty struct sent")
	}
	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Type().Field(i)
		tag := field.Tag.Get(tagName)
		//fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type, tag)

		// Skip if ignored explicitly
		if tag == "-" {
			continue
		}

		// if tag is empty or not defined check if this is a struct
		// and check for its fields inside for tags
		if tag == "" {
			if t.Field(i).Kind() == reflect.Struct {
				qVars, _ := getMapOfAllKeyValues(t.Field(i).Interface()) //recursive call
				for k, v := range *qVars {
					vars[k] = v
				}
			} else {
				continue
			}
		} else {
			// recursive check nested fields in case this is a struct
			if t.Field(i).Kind() == reflect.Struct {
				qVars, _ := getMapOfAllKeyValues(t.Field(i).Interface())
				for k, v := range *qVars {
					vars[fmt.Sprintf("%s.%s", tag, k)] = v // prepend the parent tag name
				}
			} else {
				vars[tag] = t.Field(i).Interface()
			}
		}
	}
	return &vars, nil
}

// buildMap builds the parent map and calls buildNestedMap to create the child maps based on dot notation
func buildMap(s []string, value interface{}, parent *map[string]interface{}) error {
	var obj = make(map[string]interface{})
	res := buildNestedMap(s, value, &obj)

	if parent != nil {
		if err := mergo.Merge(parent, res); err != nil {
			return err
		}
	}
	return nil
}

// ToMap creates a map based on the custom struct tag: `tag` values
// these values can be written in dot notation to create complex nested maps
// for a more comprehensive example, please see the TODO: test cases
func ToMap(obj interface{}, tag string) (*map[string]interface{}, error) {
	tagName = tag
	s, err := getMapOfAllKeyValues(obj)
	if err != nil {
		return nil, err
	}

	var parentMap = make(map[string]interface{})
	for k, v := range *s {
		keys := strings.Split(k, ".")
		if err := buildMap(keys, v, &parentMap); err != nil {
			return nil, err
		}
	}
	return &parentMap, nil
}

// buildNestedMap recursively builds a (nested) map based on dot notation
func buildNestedMap(parts []string, value interface{}, obj *map[string]interface{}) map[string]interface{} {
	if len(parts) > 1 {
		// get the first elem in list, and remove that elem from list
		var first string
		first, parts = parts[0], parts[1:]

		var m = make(map[string]interface{})
		m[first] = buildNestedMap(parts, value, obj)
		return m
	}
	return map[string]interface{}{parts[0]: value}
}
