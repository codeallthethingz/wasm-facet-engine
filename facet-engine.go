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

// CreateFacets take an json string representation of an array of objects and turn them in to facets.
// facetPaths is a query of which facets in the data to use to create facets.
func CreateFacets(jsonData string, facetPath *FacetPath) (map[string]*FacetGroup, error) {
	facetGroups := map[string]*FacetGroup{}

	if strings.TrimSpace(jsonData) == "" {
		return facetGroups, nil
	}
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
		for _, object := range arraysObject {

			name := ""
			resultingObject = object.(map[string]interface{})
			namePaths := strings.Split(facetPath.NameFieldDotNotation, ".")
			for i, path := range namePaths {
				if i == len(namePaths)-1 {
					name = getString(resultingObject, path)
					break
				}
				resultingObject = getSubObject(resultingObject, path)
			}

			nameMeta := ""
			resultingObject = object.(map[string]interface{})
			nameMetaPaths := strings.Split(facetPath.NameMetaDotNotation, ".")
			for i, path := range nameMetaPaths {
				if i == len(nameMetaPaths)-1 {
					nameMeta = getString(resultingObject, path)
					break
				}
				resultingObject = getSubObject(resultingObject, path)
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
	if thing == nil {
		return nil
	}
	return thing.(map[string]interface{})
}
func getArrays(object map[string]interface{}, path string) []interface{} {
	thing := object[path]
	if thing == nil {
		return nil
	}
	return thing.([]interface{})
}
func getString(object map[string]interface{}, path string) string {
	thing := object[path]
	if thing == nil {
		return ""
	}
	return thing.(string)
}
