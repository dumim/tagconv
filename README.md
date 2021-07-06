[![Documentation](https://godoc.org/github.com/dumim/tagconv?status.svg)](http://godoc.org/github.com/dumim/tagconv)
[![Go Report Card](https://goreportcard.com/badge/github.com/dumim/tagconv)](https://goreportcard.com/report/github.com/dumim/tagconv)
[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)

# TagConv
Convert any Go Struct to a Map based on custom struct tags with dot notation.

## Background
This package tries to simplify certain use-cases where the same struct needs to be organised/represented differently (eg: mapping data from the db to a presentable API JSON output).
This would normally have to be done by having two different structs and manually mapping the data between each other.


This package allows you to use any custom struct tag to define the mapping.
This mapping follows the dot-notation convention. Example:
```go
Hello string `mytag:"hello.world"`
```
The above will result in a map with the JSON equivalent of:
```json
{
    "hello": {
        "world": "hello world"
    }
}
```
Any number of custom tags can be used to represent the same struct in unlimited number of different ways. For examples, see below.


## Usage & Examples

Import the package
```go
import "github.com/dumim/tagconv"
```

Define your struct with custom struct tags:

```go
type MyStruct struct {
    Age   string `foo:"age"`
    Year  int    `foo:"dob.year"`
    Month int    `foo:"dob.month"`
}

obj := MyStruct{
    Age:   "22",
    Year:  1998,
    Month: 1,
}

tagName = "foo"
myMap, err := ToMap(obj, tagName)
if err != nil {
    panic()
}
```
This will result in a map that looks like:
```go
myMap = map[string]interface{}{
	"age": "22",
	"dob": map[string]interface{}{
	    "year": 1998,
	    "month": 1,
    }
}
```
Converting to JSON ...
```go
myMapJSON, err := json.MarshalIndent(myMap, "", "    ")
    if err != nil {
    panic()
}
fmt.Print(myMapJSON)
```
... will result in something similar to:
```json
{
  "age": "22",
  "dob": {
    "year": 1998,
    "month": 1
  }
}
```

---
### Multiple struct tags

You can use multiple struct tags for different representation of the same struct:
For example, similar to the previous example:
```go
type MyStructMultiple struct {
    Age   string `foo:"age" bar:"details.my_age"`
}

obj := MyStruct{
    Age:   "22",
}
```
Using `tagconv` for `obj` over the `foo` tag (`ToMap(obj, "foo")`) will result in:
```json
{
  "age": "22"
}
```
whereas using `bar` (`ToMap(obj, "bar")`) on the same `obj` will result in:
```json
{
  "details": {
    "my_age": "22"
  }
}
```
---
### Tag options
- If a nested struct has a tag, this will create a parent-child relationship with that tag and the tags of the fields within that struct.
- Dot notation will create a parent-child relationship for every `.`.
- Not setting any tag will ignore that field, unless if it's a struct; then it will go inside the struct to check its tags
- `-` will explicitly ignore that field. As opposed to above, it will not look inside even if the field is of struct type.

For an example that includes all the above scenarios see the code below:


### More complex example

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
---

## Testing
Run the go tests using `go test ./.. -v`


---

## Acknowledgements

- [Helpful Stackoverflow answer](https://stackoverflow.com/a/7794127/10340220)
- [Mergo](https://github.com/imdario/mergo)
## Contributing

Contributions are always welcome!
