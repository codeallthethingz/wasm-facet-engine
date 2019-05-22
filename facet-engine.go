package engine

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FacetGroup contains the description of a facet.
type FacetGroup struct {
	Name   string
	Facets map[string]*Facet
}

// Facet contains the values of a facet
type Facet struct {
	Name   string
	Values *Set
}

// FacetPath How to get out data from
type FacetPath struct {
	ArrayDotNotation     string
	NameMetaDotNotation  string
	NameFieldDotNotation string
	ValueMapDotNotation  string
}

func getAtPath(data map[string]interface{}, path []string) interface{} {
	if len(path) == 1 {
		return data[path[0]]
	}
	return getAtPath(data[path[0]].(map[string]interface{}), path[1:])
}

// CreateFacets take an json string representation of an array of objects and turn them in to facets.
// facetPaths is a query of which facets in the data to use to create facets.
func CreateFacets(jsonData string, facetPath *FacetPath) (map[string]*FacetGroup, error) {
	facetGroups := map[string]*FacetGroup{}

	if strings.TrimSpace(jsonData) == "" {
		return facetGroups, nil
	}
	var genericObjects []map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &genericObjects)
	if err != nil {
		return nil, err
	}

	for _, genericObject := range genericObjects {

		arrayPaths := strings.Split(facetPath.ArrayDotNotation, ".")
		namePaths := strings.Split(facetPath.NameFieldDotNotation, ".")
		nameMetaPaths := strings.Split(facetPath.NameMetaDotNotation, ".")
		arraysObject := getAtPath(genericObject, arrayPaths).([]interface{})
		for _, object := range arraysObject {
			o := object.(map[string]interface{})
			name := getAtPath(o, namePaths).(string)
			nameMeta := getAtPath(o, nameMetaPaths).(string)
			valuePaths := strings.Split(facetPath.ValueMapDotNotation, ".")
			values := toStringMap(getAtPath(o, valuePaths).(map[string]interface{}))
			key := strings.ToLower(fmt.Sprintf("%s (%s)", name, nameMeta))
			fmt.Println(key)
			if len(values) == 0 || strings.TrimSpace(name) == "" || strings.TrimSpace(nameMeta) == "" {
				continue
			}
			if _, ok := facetGroups[key]; !ok {
				fmt.Println("creating new facetGroup")
				facetGroups[key] = &FacetGroup{
					Name:   key,
					Facets: map[string]*Facet{},
				}
			} else {
				fmt.Println("facet group already exist")
			}

			for k, v := range values {
				facetKey := strings.ToLower(k)
				facet, ok := facetGroups[key].Facets[facetKey]
				if ok {
					fmt.Printf("adding new value: %v, %s, %s\n", facet.Values, facetKey, v)
				} else {
					fmt.Printf("creating new value: %s, %s\n", facetKey, v)
					facetGroups[key].Facets[facetKey] = &Facet{
						Name:   facetKey,
						Values: NewSet(),
					}
				}
				facetGroups[key].Facets[facetKey].Values.Add(v)
			}
		}
	}

	return facetGroups, nil
}
