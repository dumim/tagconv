package tagconv

import (
	"dario.cat/mergo"
	"fmt"
	"reflect"
	"strings"
)

var tagName = "" // initialise the struct tag value

const omitEmptyTagOption = "omitempty" // omit values with this tag option if empty

// getMapOfAllKeyValues builds a map of the fully specified key and the value from the struct tag
// the struct tags with the full dot notation will be used as the key, and the value as the value
// slices will be also be maps
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
	"list":{
			{"name":"hi", "value":1},
			{"name":"world", "value":2}
		}
*/
func getMapOfAllKeyValues(s interface{}) *map[string]interface{} {
	var vars = make(map[string]interface{}) // this will hold the variables as a map (JSON)

	// get value of object
	t := reflect.ValueOf(s)
	if t.IsZero() {
		return nil
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
				// only check if the value can be obtained without panicking (eg: for unexported fields)
				if t.Field(i).CanInterface() {
					qVars := getMapOfAllKeyValues(t.Field(i).Interface()) //recursive call
					if qVars != nil {
						for k, v := range *qVars {
							vars[k] = v
						}
					}
				}
			} else {
				continue
			}
		} else {
			// omitempty tag passed?
			tag, shouldOmitEmpty := shouldOmitEmpty(tag) // overwrite tag
			// recursive check nested fields in case this is a struct
			if t.Field(i).Kind() == reflect.Struct {
				// only check if the value can be obtained without panicking (eg: for unexported fields)
				if t.Field(i).CanInterface() {
					if shouldOmitEmpty {
						if t.Field(i).IsZero() {
							continue
						}
					}
					qVars := getMapOfAllKeyValues(t.Field(i).Interface()) //recursive call
					if qVars != nil {
						for k, v := range *qVars {
							vars[fmt.Sprintf("%s.%s", tag, k)] = v // prepend the parent tag name
						}
					}
				}
			} else {
				// only check if the value can be obtained without panicking (eg: for unexported fields)
				if t.Field(i).CanInterface() {
					if shouldOmitEmpty {
						if t.Field(i).IsZero() {
							continue
						}
					}
					vars[tag] = t.Field(i).Interface()
				}
			}
		}
	}

	// process slices separately
	// and create the final map
	var finalMap = make(map[string]interface{})
	// iterate through the map
	for k, v := range vars {
		if reflect.TypeOf(v) != nil {
			switch reflect.TypeOf(v).Kind() {
			// if any of them is a slice
			case reflect.Slice:
				if reflect.TypeOf(v).Elem().Kind() == reflect.Struct {
					var sliceOfMap []map[string]interface{}
					s := reflect.ValueOf(v)
					// iterate through the slice
					for i := 0; i < s.Len(); i++ {
						if s.Index(i).CanInterface() {
							m := getMapOfAllKeyValues(s.Index(i).Interface()) // get the map value of the object, recursively
							if m != nil {
								sliceOfMap = append(sliceOfMap, *m) // append to the slice
							}
						}
					}
					finalMap[k] = sliceOfMap
				} else {
					finalMap[k] = v
				}
			default:
				finalMap[k] = v
			}
		}
	}

	return &finalMap
}

// shouldOmitEmpty checks if the omitEmptyTagOption option is passed in the tag
// eg: `foo:"bar,omitempty"`
func shouldOmitEmpty(originalTag string) (string, bool) {
	if ss := strings.Split(originalTag, ","); len(ss) > 1 {
		// TODO: add more validation & error checking
		if strings.TrimSpace(ss[1]) == omitEmptyTagOption {
			return ss[0], true
		}
		return ss[0], false
	} else {
		return originalTag, false
	}
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
// for a more comprehensive example, please see the
func ToMap(obj interface{}, tag string) (*map[string]interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("nil object passed")
	}

	tagName = tag
	s := getMapOfAllKeyValues(obj)

	if s == nil {
		return nil, fmt.Errorf("no valid map could be formed")
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
