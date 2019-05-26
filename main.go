package main

import (
	"encoding/json"
	"fmt"

	"github.com/gopherjs/gopherwasm/js"
)

var facetGroups *FacetGroups
var query = &Query{}

// JSQuery wasm interface to query the facet groups
func JSQuery(args []js.Value) {
	fmt.Println("query called")
	ids, err := facetGroups.Query(query)
	if err != nil {
		panic(err)
	}
	idsByets, err := json.Marshal(ids)
	if err != nil {
		panic(err)
	}

	js.Global().Call("facetEngineCallbackRecords", string(idsByets))
}

// JSInitializeObjects wasm interface to take the data and parse out the facets
func JSInitializeObjects(args []js.Value) {
	fmt.Println("facetEngineInitializeObjects called")
	configString := args[0].String()
	dataJSON := args[1].String()

	fmt.Println("Config: " + configString)
	fmt.Println("Data: " + dataJSON)

	facetPath := &FacetPath{}
	err := json.Unmarshal([]byte(configString), facetPath)
	if err != nil {
		panic(err)
	}
	facetGroups, err = CreateFacetGroups(dataJSON, facetPath)
	if err != nil {
		panic(err)
	}
	facetGroupsBytes, err := json.Marshal(facetGroups)
	if err != nil {
		panic(err)
	}
	fmt.Println("facets: " + string(facetGroupsBytes))
	js.Global().Call("facetEngineCallbackFacets", string(facetGroupsBytes))
}

func registerCallbacks() {
	js.Global().Set("facetEngineInitializeObjects", js.NewCallback(JSInitializeObjects))
	js.Global().Set("facetEngineQuery", js.NewCallback(JSQuery))
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("hello")
	registerCallbacks()
	<-c
}
