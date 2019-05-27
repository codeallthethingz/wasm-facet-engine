package main

import (
	"encoding/json"
	"fmt"

	"github.com/gopherjs/gopherwasm/js"
)

var facetEngine *FacetEngine

func main() {
	// create empty channel so main doesn't exit when it's wasm-ed.
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}

// JSClearFilters remove all the filters
func JSClearFilters(args []js.Value) {
	fmt.Println("clear filters")
	facetEngine.ClearFilters()
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
	addFilter(facetGroupName, facetName, inclusiveMin, min, inclusiveMax, max)
}

func addFilter(facetGroupName string, facetName string, inclusiveMin bool, min float64, inclusiveMax bool, max float64) error {
	minRange := Exclusive(min)
	maxRange := Exclusive(max)
	if inclusiveMin {
		minRange = Inclusive(min)
	}
	if inclusiveMax {
		maxRange = Inclusive(max)
	}
	return facetEngine.AddFilter(facetGroupName, facetName, minRange, maxRange)
}

// JSQuery wasm interface to query the facet groups
func JSQuery(args []js.Value) {
	fmt.Println("query called")
	ids, facetGroups, err := query()
	if err != nil {
		panic(err)
	}
	js.Global().Call("facetEngineCallbackRecords", ids)
	js.Global().Call("facetEngineCallbackFacets", facetGroups)
}

func query() (string, string, error) {
	ids, facetGroups, err := facetEngine.Query()
	if err != nil {
		return "", "", err
	}
	idsByets, err := json.Marshal(ids)
	if err != nil {
		return "", "", err
	}
	facetGroupBytes, err := json.Marshal(facetGroups)
	if err != nil {
		return "", "", err
	}
	return string(idsByets), string(facetGroupBytes), nil
}

// JSInitializeObjects wasm interface to take the data and parse out the facets
func JSInitializeObjects(args []js.Value) {
	fmt.Println("facetEngineInitializeObjects called")
	configString := args[0].String()
	dataJSON := args[1].String()
	facetGroupsString, err := initializeObjects(configString, dataJSON)
	if err != nil {
		panic(err)
	}
	js.Global().Call("facetEngineCallbackFacets", facetGroupsString)
}

func initializeObjects(configString string, dataJSON string) (string, error) {
	fmt.Println("Config: " + configString)
	fmt.Println("Data: " + dataJSON)

	facetPath := &FacetPath{}
	err := json.Unmarshal([]byte(configString), facetPath)
	if err != nil {
		return "", err
	}
	var facetGroups map[string]*FacetGroup
	facetEngine, facetGroups, err = NewFacetEngine(dataJSON, facetPath)
	if err != nil {
		return "", err
	}
	facetGroupsBytes, err := json.Marshal(facetGroups)
	if err != nil {
		return "", err
	}
	fmt.Println("facets: " + string(facetGroupsBytes))
	return string(facetGroupsBytes), nil
}

func registerCallbacks() {
	js.Global().Set("facetEngineInitializeObjects", js.NewCallback(JSInitializeObjects))
	js.Global().Set("facetEngineQuery", js.NewCallback(JSQuery))
	js.Global().Set("facetEngineAddFilter", js.NewCallback(JSAddFilter))
	js.Global().Set("facetEngineClearFilters", js.NewCallback(JSClearFilters))
}
