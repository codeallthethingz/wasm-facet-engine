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
	IDDotNotation        string
	ArrayDotNotation     string
	NameMetaDotNotation  string
	NameFieldDotNotation string
	ValueMapDotNotation  string
}

// CreateFacetGroups take an json string representation of an array of objects and turn them in to facets.
// facetPaths is a query of which facets in the data to use to create facets.
func CreateFacetGroups(jsonData string, facetPath *FacetPath) (map[string]*FacetGroup, error) {
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
		idPaths := facetPath.IDDotNotation
		if idPaths == "" {
			idPaths = "id"
		}
		id := getAtPathString(genericObject, strings.Split(idPaths, "."))
		if strings.TrimSpace(id) == "" {
			return nil, fmt.Errorf("found record with no id")
		}
		arrayPaths := strings.Split(facetPath.ArrayDotNotation, ".")
		namePaths := strings.Split(facetPath.NameFieldDotNotation, ".")
		nameMetaPaths := strings.Split(facetPath.NameMetaDotNotation, ".")
		valuePaths := strings.Split(facetPath.ValueMapDotNotation, ".")
		arraysObject := getAtPathArray(genericObject, arrayPaths)
		for _, object := range arraysObject {
			o := object.(map[string]interface{})
			name := getAtPathString(o, namePaths)
			nameMeta := getAtPathString(o, nameMetaPaths)
			values := getAtPathMap(o, valuePaths)
			key := strings.ToLower(fmt.Sprintf("%s (%s)", name, nameMeta))
			if len(values) == 0 || strings.TrimSpace(name) == "" || strings.TrimSpace(nameMeta) == "" {
				continue
			}
			if _, ok := facetGroups[key]; !ok {
				facetGroups[key] = &FacetGroup{
					Name:   key,
					Facets: map[string]*Facet{},
				}
			}

			for k, v := range values {
				facetKey := strings.ToLower(k)
				if _, ok := facetGroups[key].Facets[facetKey]; !ok {
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

func getAtPathArray(data map[string]interface{}, path []string) []interface{} {
	obj := getAtPath(data, path)
	if obj == nil {
		return nil
	}
	return obj.([]interface{})
}

func getAtPathString(data map[string]interface{}, path []string) string {
	obj := getAtPath(data, path)
	if obj == nil {
		return ""
	}
	return obj.(string)
}
func getAtPathMap(data map[string]interface{}, path []string) map[string]string {
	obj := getAtPath(data, path)
	if obj == nil {
		return nil
	}
	values := map[string]string{}
	for key, value := range obj.(map[string]interface{}) {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		values[strKey] = strValue
	}
	return values
}

func getAtPath(data map[string]interface{}, path []string) interface{} {
	if len(path) == 1 {
		return data[path[0]]
	}
	if data[path[0]] == nil {
		return nil
	}
	return getAtPath(data[path[0]].(map[string]interface{}), path[1:])
}
