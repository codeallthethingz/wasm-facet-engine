package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// FacetGroups is a map of FacetGroup
type FacetGroups struct {
	FacetGroup   map[string]*FacetGroup
	RecordLookup RecordLookup
}

// RecordLookup Set of records
type RecordLookup map[string][]*Record

// Add a record to the map.
func (r RecordLookup) Add(key string, record *Record) {
	if _, ok := r[key]; !ok {
		r[key] = []*Record{}
	}
	r[key] = append(r[key], record)
}

// Record holds values and ids for filtering records.
type Record struct {
	Value string
	ID    string
}

// NewFacetGroups create a new one.
func NewFacetGroups() *FacetGroups {
	return &FacetGroups{
		FacetGroup:   map[string]*FacetGroup{},
		RecordLookup: map[string][]*Record{},
	}
}

// Get return a facetgroup by key
func (f *FacetGroups) Get(key string) *FacetGroup {
	if value, ok := f.FacetGroup[key]; ok {
		return value
	}
	return nil
}

// Lookup return a facetgroup by key and a boolean of whether it exists
func (f *FacetGroups) Lookup(key string) (*FacetGroup, bool) {
	value, ok := f.FacetGroup[key]
	return value, ok
}

// Len number of facet groups
func (f *FacetGroups) Len() int {
	return len(f.FacetGroup)
}

// Set set a facetgroup by key
func (f *FacetGroups) Set(key string, value *FacetGroup) {
	f.FacetGroup[key] = value
}

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

// Query represents a set of filters to be applied to the data.
type Query struct {
	Filters []filter
}

type filter struct {
	FacetGroupName string
	FacetName      string
	Min            Range
	Max            Range
}

// AddFilter adds a set of criteria that records will have to match.
func (q *Query) AddFilter(facetGroupName string, facetName string, min Range, max Range) *Query {
	if q.Filters == nil {
		q.Filters = []filter{}
	}
	q.Filters = append(q.Filters, filter{
		FacetGroupName: facetGroupName,
		FacetName:      facetName,
		Min:            min,
		Max:            max,
	})
	return q
}

// Range repnesents min and max bounds inclusive or exclusive
type Range interface {
	IsInclusive() bool
	Value() float64
}

// Inclusive range value
func Inclusive(value float64) Range {
	return inclusive{
		value: value,
	}
}

// Exclusive range value
func Exclusive(value float64) Range {
	return exclusive{
		value: value,
	}
}

type inclusive struct {
	value float64
}

func (i inclusive) IsInclusive() bool {
	return true
}
func (i inclusive) Value() float64 {
	return i.value
}

type exclusive struct {
	value float64
}

func (e exclusive) IsInclusive() bool {
	return false
}
func (e exclusive) Value() float64 {
	return e.value
}

// Query filter the records and return ids that match the filters
func (f FacetGroups) Query(query *Query) ([]string, error) {
	listOfMaps := make([]map[string]bool, len(query.Filters))
	results := []string{}
	for i, filter := range query.Filters {
		records, ok := f.RecordLookup[fmt.Sprintf("%s - %s", filter.FacetGroupName, filter.FacetName)]
		if ok {
			matchingIDMap, err := toStringMap(records, filter)
			if err != nil {
				return nil, err
			}
			listOfMaps[i] = matchingIDMap
		}
	}
	for k := range listOfMaps[0] {
		inAll := true
		for i := 1; i < len(listOfMaps); i++ {
			_, inThis := listOfMaps[i][k]
			inAll = inThis && inAll
		}
		if inAll {
			results = append(results, k)
		}
	}
	return results, nil
}

func toStringMap(records []*Record, filter filter) (map[string]bool, error) {
	results := map[string]bool{}
	for _, record := range records {
		value, err := strconv.ParseFloat(record.Value, 64)
		if err != nil {
			return nil, err
		}
		if ((value >= filter.Min.Value() && filter.Min.IsInclusive()) || (value > filter.Min.Value() && !filter.Min.IsInclusive())) &&
			((value <= filter.Max.Value() && filter.Max.IsInclusive()) || (value < filter.Max.Value() && !filter.Max.IsInclusive())) {
			results[record.ID] = true
		}
	}
	return results, nil
}

// CreateFacetGroups take an json string representation of an array of objects and turn them in to facets.
// facetPaths is a query of which facets in the data to use to create facets.
func CreateFacetGroups(jsonData string, facetPath *FacetPath) (*FacetGroups, error) {
	facetGroups := NewFacetGroups()

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
			if _, ok := facetGroups.Lookup(key); !ok {
				facetGroups.Set(key, &FacetGroup{
					Name:   key,
					Facets: map[string]*Facet{},
				})
			}

			for k, v := range values {
				facetKey := strings.ToLower(k)
				facetGroups.RecordLookup.Add(fmt.Sprintf("%s - %s", key, facetKey), &Record{
					ID:    id,
					Value: v,
				})
				facetGroup := facetGroups.Get(key)
				if _, ok := facetGroup.Facets[facetKey]; !ok {
					facetGroup.Facets[facetKey] = &Facet{
						Name:   facetKey,
						Values: NewSet(),
					}
				}
				facetGroup.Facets[facetKey].Values.Add(v)
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
