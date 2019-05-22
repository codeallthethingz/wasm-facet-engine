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
	genericObjects := []map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonData), &genericObjects)
	if err != nil {
		return nil, err
	}

	for _, genericObject := range genericObjects {

		arrayPaths := strings.Split(facetPath.ArrayDotNotation, ".")
		resultingObject := genericObject
		arraysObject := []interface{}{}
		for i, path := range arrayPaths {
			if i == len(arrayPaths)-1 {
				arraysObject = getArrays(resultingObject, path)
				break
			}
			resultingObject = getSubObject(resultingObject, path)
		}

		namePaths := strings.Split(facetPath.NameFieldDotNotation, ".")
		nameMetaPaths := strings.Split(facetPath.NameMetaDotNotation, ".")
		for _, object := range arraysObject {
			o := object.(map[string]interface{})
			name := getAtPath(o, namePaths).(string)
			nameMeta := getAtPath(o, nameMetaPaths).(string)

			resultingObject = object.(map[string]interface{})

			for i, path := range nameMetaPaths {
				if i == len(nameMetaPaths)-1 {
					nameMeta = getString(resultingObject, path)
					break
				}
				resultingObject = getSubObject(resultingObject, path)
			}

			key := fmt.Sprintf("%s (%s)", name, nameMeta)
			fmt.Println(key)
			if _, ok := facetGroups[key]; !ok {
				fmt.Println("creating new facetGroup")
				facetGroups[key] = &FacetGroup{
					Name:   key,
					Facets: map[string]*Facet{},
				}
			} else {
				fmt.Println("facet group already exist")
			}

			values := map[string]string{}
			resultingObject = object.(map[string]interface{})
			valuePaths := strings.Split(facetPath.ValueMapDotNotation, ".")
			for i, path := range valuePaths {
				if i == len(valuePaths)-1 {
					values = getMap(resultingObject, path)
					break
				}
				resultingObject = getSubObject(resultingObject, path)
			}
			for k, v := range values {
				facet, ok := facetGroups[key].Facets[k]
				if ok {
					fmt.Printf("adding new value: %v, %s, %s\n", facet.Values, k, v)
				} else {
					fmt.Printf("creating new value: %s, %s\n", k, v)
					facetGroups[key].Facets[k] = &Facet{
						Name:   k,
						Values: NewSet(),
					}
				}
				facetGroups[key].Facets[k].Values.Add(v)
			}
		}
	}

	return facetGroups, nil
}

func getMap(object map[string]interface{}, path string) map[string]string {
	thing := object[path]
	values := map[string]string{}
	for key, value := range thing.(map[string]interface{}) {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		values[strKey] = strValue
	}
	return values
}

func getSubObject(object map[string]interface{}, path string) map[string]interface{} {
	thing := object[path]
	return thing.(map[string]interface{})
}
func getArrays(object map[string]interface{}, path string) []interface{} {
	thing := object[path]
	return thing.([]interface{})
}
func getString(object map[string]interface{}, path string) string {
	thing := object[path]
	return thing.(string)
}
