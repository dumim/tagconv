package tagconv

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

// structs with the custom tags on
type Obj struct {
	Name  string `custom:"name"`
	Text  string `custom:"text"`
	World string `custom:"data.world"`
}
type ObjTwo struct {
	Hello string `custom:"hello"`
	Text  string `custom:"data.text"`
}
type ObjThree struct {
	Name  string `custom:"name"`
	Value int    `custom:"value"`
}
type Example struct {
	Name     string     `custom:"name"`
	Email    string     `custom:"email"`
	Obj      Obj        `custom:"object"`
	ObjTwo   ObjTwo     // no tag
	ObjThree ObjTwo     `custom:"-"` // explicitly ignored
	Id       int        `custom:"id"`
	Call     int        `custom:"data.call"`
	ArrayObj []ObjThree `custom:"list"`
	//three    int    `custom:"three"` // unexported, TODO: handle panic
}

// TestFullStructToMap calls ToMap function and checks against the expected result
// The struct used tries to cover all the scenarios
func TestFullStructToMap(t *testing.T) {

	// the initial object
	initial := Example{
		Name:  "2",
		Email: "3",
		Obj: Obj{
			Name:  "4",
			Text:  "5",
			World: "6",
		},
		ObjTwo: ObjTwo{
			Hello: "1",
			Text:  "2",
		},
		Id:   01,
		Call: 02,
		ArrayObj: []ObjThree{
			{"hi", 1},
			{"world", 2},
		},
	}

	// expected response
	expectedJSON := `{
			"name": "2",
			"email": "3",
			"object": {
			  "name": "4",
			  "text": "5",
			  "data": {
				"world": "6"
			  }
			},
			"hello": "1",
			"data": {
			  "text": "2",
			  "call": 2
			},
			"id": 1,
			"list": [
				{
					"name": "hi",
					"value": 1
				},
				{
					"name": "world",
					"value": 2
				}
			]
		  }
		`

	// get the map from custom tags
	tagName = "custom"
	actual, err := ToMap(initial, tagName)
	if err != nil {
		t.Fail()
	}

	// convert to json to compare
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fail()
	}

	// compare
	require.JSONEqf(t, expectedJSON, string(actualJSON), "JSON mismatch")
}

// TestMultipleTagsStructToMap calls ToMap function and checks against the expected result
// The struct used tries to use multiple struct tags for different responses
func TestMultipleTagsStructToMap(t *testing.T) {
	type MyStruct struct {
		Age   string `foo:"age" bar:"details.myAge"`
		Year  int    `foo:"dob.year" bar:"details.birthYear"`
		Month int    `foo:"dob.month" bar:"-"`
	}

	obj := MyStruct{
		Age:   "22",
		Year:  1998,
		Month: 1,
	}

	// expected response
	expectedJSONOne := `{
		  "age": "22",
		  "dob": {
			"year": 1998,
			"month": 1
		  }
		}
	`
	expectedJSONTwo := `{
		  "details": {
			"myAge": "22",
			"birthYear": 1998
		  }
		}
	`

	// get the map from custom tags for tag "foo"
	actualOne, err := ToMap(obj, "foo")
	if err != nil {
		t.Fail()
	}
	actualJSONOne, err := json.Marshal(actualOne)
	if err != nil {
		t.Fail()
	}

	// get the map from custom tags for tag "bar"
	actualTwo, err := ToMap(obj, "bar")
	if err != nil {
		t.Fail()
	}
	actualJSONTwo, err := json.Marshal(actualTwo)
	if err != nil {
		t.Fail()
	}

	// compare
	require.JSONEqf(t, expectedJSONOne, string(actualJSONOne), "JSON mismatch for foo tags")
	require.JSONEqf(t, expectedJSONTwo, string(actualJSONTwo), "JSON mismatch for bar tags")
}
