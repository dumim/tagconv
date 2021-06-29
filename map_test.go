package tagconv

import (
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

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func _TestHelloName(t *testing.T) {

	// the initial object
	_ := Example{
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
	_ = `{
			"name": "",
			"email": "",
			"object": {
			  "name": "",
			  "text": "",
			  "data": {
				"world": ""
			  }
			},
			"hello": "",
			"data": {
			  "text": "",
			  "call": 0
			},
			"id": 0,
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
	// TODO: do the testing
}

func TestPass(t *testing.T) {
}
