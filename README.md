[![Documentation](https://godoc.org/github.com/dumim/tagconv?status.svg)](http://godoc.org/github.com/dumim/tagconv)
[![Go Report Card](https://goreportcard.com/badge/github.com/dumim/tagconv)](https://goreportcard.com/report/github.com/dumim/tagconv)
[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)

# TagConv
Convert any Go Struct to a Map based on custom struct tags with dot notation

TODO: bacgkround
## Usage/Examples

Import the package
```go
import "github.com/dumim/tagconv"
```

Given a deeply-nested complex struct with custom tags like below:
```go
type Obj struct {
	Name  string `custom:"name"`
	Text  string `custom:"text"`
	World string `custom:"data.world"` // dot notation inside nested struct
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
	Call     int        `custom:"data.call"` // top-level dot notation
	ArrayObj []ObjThree `custom:"list"`
}
```
The `ToMap` function can be used to convert this into a JSON/Map based on the values defined in the given custom tag like so.
```go
obj := Example{
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

// get the map from custom tags
tagName = "custom"
myMap, err := ToMap(obj, tagName)
if err != nil {
    panic()
}

myMapJSON, err := json.MarshalIndent(myMap, "", "    ")
if err != nil {
    panic()
}

fmt.Print(myMapJSON)

```
This will produce a result similar to:
```json
{
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
```
  
