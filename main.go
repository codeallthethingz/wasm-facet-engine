package main

import (
	"encoding/json"
	"fmt"

	"github.com/gopherjs/gopherwasm/js"
)

var facetEngine *FacetEngine
var query = &Query{}

// JSClearFilters remove all the filters
func JSClearFilters(args []js.Value) {
	fmt.Println("clear filters")
	query.Filters = []filter{}
}

// JSAddFilter addes a filter to the query object
func JSAddFilter(args []js.Value) {
	fmt.Println("add filter called")
	facetGroupName := args[0].String()
	facetName := args[1].String()
	inclusiveMin := args[2].Bool()
	min := args[3].Float()
	inclusiveMax := args[4].Bool()
	max := args[5].Float()
	minRange := Exclusive(min)
	maxRange := Exclusive(max)
	if inclusiveMin {
		minRange = Inclusive(min)
	}
	if inclusiveMax {
		maxRange = Inclusive(max)
	}
	query.AddFilter(facetGroupName, facetName, minRange, maxRange)
}

// JSQuery wasm interface to query the facet groups
func JSQuery(args []js.Value) {
	fmt.Println("query called")
	ids, facetGroups, err := facetEngine.Query(query)
	if err != nil {
		panic(err)
	}
	idsByets, err := json.Marshal(ids)
	if err != nil {
		panic(err)
	}
	facetGroupBytes, err := json.Marshal(facetGroups)
	if err != nil {
		panic(err)
	}

	js.Global().Call("facetEngineCallbackRecords", string(idsByets))
	js.Global().Call("facetEngineCallbackFacets", string(facetGroupBytes))
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
	var facetGroups map[string]*FacetGroup
	facetEngine, facetGroups, err = NewFacetEngine(dataJSON, facetPath)
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
	js.Global().Set("facetEngineAddFilter", js.NewCallback(JSAddFilter))
	js.Global().Set("facetEngineClearFilters", js.NewCallback(JSClearFilters))
}

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}
