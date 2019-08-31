package main

import (
	"encoding/json"
	"github.com/gopherjs/gopherwasm/js"
)

var facetEngine *FacetEngine

func main() {
	// create empty channel so main doesn't exit when it's wasm-ed.
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	js.Global().Get("facetEngine").Set("initializeObjects", js.NewCallback(JSInitializeObjects))
	js.Global().Get("facetEngine").Set("query", js.NewCallback(JSQuery))
	js.Global().Get("facetEngine").Set("addFilter", js.NewCallback(JSAddFilter))
	js.Global().Get("facetEngine").Set("clearFilters", js.NewCallback(JSClearFilters))
}

// JSClearFilters remove all the filters
//noinspection GoUnusedParameter
func JSClearFilters(args []js.Value) {
	facetEngine.ClearFilters()
}

// JSAddFilter adds a filter to the query object
func JSAddFilter(args []js.Value) {
	facetGroupName := args[0].String()
	facetName := args[1].String()
	inclusiveMin := args[2].Bool()
	min := args[3].Float()
	inclusiveMax := args[4].Bool()
	max := args[5].Float()
	err := addFilter(facetGroupName, facetName, inclusiveMin, min, inclusiveMax, max)
	if err != nil {
		panic(err)
	}
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

// JSQuery WASM interface to query the facet groups
func JSQuery(args []js.Value) {
	ids, facetGroups, err := query()
	if err != nil {
		panic(err)
	}
	args[0].Invoke(ids)
	args[1].Invoke(facetGroups)
}

func query() (string, string, error) {
	ids, facetGroups, err := facetEngine.Query()
	if err != nil {
		return "", "", err
	}
	idsBytes, err := json.Marshal(ids)
	if err != nil {
		return "", "", err
	}
	facetGroupBytes, err := json.Marshal(facetGroups)
	if err != nil {
		return "", "", err
	}
	return string(idsBytes), string(facetGroupBytes), nil
}

// JSInitializeObjects wasm interface to take the data and parse out the facets
func JSInitializeObjects(args []js.Value) {
	configString := args[0].String()
	dataJSON := args[1].String()
	facetGroupsString, err := initializeObjects(configString, dataJSON)
	if err != nil {
		panic(err)
	}
	args[2].Invoke(facetGroupsString)
}

func initializeObjects(configString string, dataJSON string) (string, error) {
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
	return string(facetGroupsBytes), nil
}
