package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FacetEngine is a map of FacetGroup
type FacetEngine struct {
	RecordLookup   RecordLookup
	facetPath      *FacetPath
	ids            *Set
	allIds         *Set
	query          *Query
	initialized    bool
	genericObjects []map[string]interface{}
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

// NewFacetEngine create a new one.
func NewFacetEngine(dataJSON string, config *FacetPath) (*FacetEngine, map[string]*FacetGroup, error) {
	facetEngine := &FacetEngine{
		RecordLookup: map[string][]*Record{},
		ids:          NewSet(),
		allIds:       NewSet(),
		query:        &Query{},
	}
	facetGroups, err := facetEngine.Initialize(dataJSON, config)
	return facetEngine, facetGroups, err
}

// FacetGroup contains the description of a facet.
type FacetGroup struct {
	Name   string            `json:"name,omitempty"`
	Facets map[string]*Facet `json:"facets,omitempty"`
}

// Facet contains the values of a facet
type Facet struct {
	Name   string `json:"name,omitempty"`
	Values *Set   `json:"values,omitempty"`
}

// FacetPath How to get out data from
type FacetPath struct {
	IDDotNotation        string `json:"idDotNotation,omitempty"`
	ArrayDotNotation     string `json:"arrayDotNotation,omitempty"`
	NameMetaDotNotation  string `json:"nameMetaDotNotation,omitempty"`
	NameFieldDotNotation string `json:"nameFieldDotNotation,omitempty"`
	ValueMapDotNotation  string `json:"valueMapDotNotation,omitempty"`
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
func (f *FacetEngine) AddFilter(facetGroupName string, facetName string, min Range, max Range) error {
	if f.query.Filters == nil {
		f.query.Filters = []filter{}
	}
	if strings.TrimSpace(facetGroupName) == "" {
		return fmt.Errorf("must specify facetgroup name")
	}
	if strings.TrimSpace(facetName) == "" {
		return fmt.Errorf("must specify facet name")
	}
	f.query.Filters = append(f.query.Filters, filter{
		FacetGroupName: facetGroupName,
		FacetName:      facetName,
		Min:            min,
		Max:            max,
	})
	return nil
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

// ClearFilters remove all the query state
func (f *FacetEngine) ClearFilters() {
	f.query = &Query{}
	f.resetAllIds()
}

func (f *FacetEngine) resetAllIds() {
	f.ids = NewSet()
	for _, id := range f.allIds.ToArray() {
		f.ids.Add(id)
	}
}

// Query filter the records and return ids that match the filters
func (f FacetEngine) Query() ([]string, map[string]*FacetGroup, error) {
	if len(f.query.Filters) == 0 {
		fmt.Println("no filters returning all")
		facetGroups, err := f.GetFacets()
		f.resetAllIds()
		return f.allIds.ToArray(), facetGroups, err
	}
	if f.ids.Len() == 0 {
		fmt.Println("results restricted to nothing, returning nothing.")
		return []string{}, map[string]*FacetGroup{}, nil
	}
	then := now()
	listOfMaps := make([]map[string]bool, len(f.query.Filters))
	f.ids = NewSet()
	for i, filter := range f.query.Filters {
		key := fmt.Sprintf("%s - %s", filter.FacetGroupName, filter.FacetName)
		if records, ok := f.RecordLookup[key]; ok {
			listOfMaps[i] = toStringMap(records, filter)
		}
	}
	for k := range listOfMaps[0] {
		inAll := true
		for i := 1; i < len(listOfMaps); i++ {
			_, inThis := listOfMaps[i][k]
			// todo should just break when the first false happens
			inAll = inThis && inAll
		}
		if inAll {
			f.ids.Add(k)
		}
	}
	fmt.Printf("search time: %d\n", now()-then)
	facetGroups, err := f.GetFacets()
	return f.ids.ToArray(), facetGroups, err
}

func now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func toStringMap(records []*Record, filter filter) map[string]bool {
	results := map[string]bool{}
	for _, record := range records {
		// this parse error is guaranteed not to happen elsewhere.
		value, _ := strconv.ParseFloat(record.Value, 64)
		if ((value >= filter.Min.Value() && filter.Min.IsInclusive()) || (value > filter.Min.Value() && !filter.Min.IsInclusive())) &&
			((value <= filter.Max.Value() && filter.Max.IsInclusive()) || (value < filter.Max.Value() && !filter.Max.IsInclusive())) {
			results[record.ID] = true
		}
	}
	return results
}

// Initialize take an json string representation of an array of objects and turn them in to facets.
// facetPaths is a query of which facets in the data to use to create facets.
func (f *FacetEngine) Initialize(jsonData string, facetPath *FacetPath) (map[string]*FacetGroup, error) {
	if strings.TrimSpace(jsonData) == "" {
		return f.GetFacets()
	}
	var genericObjects []map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &genericObjects)
	if err != nil {
		return nil, err
	}
	f.genericObjects = genericObjects
	f.facetPath = facetPath

	facetGroups, err := f.GetFacets()
	f.resetAllIds()
	f.initialized = true
	return facetGroups, err
}

// GetFacets return a list of facets for the list of ids.  If ids is nil, return all possible facets.
func (f *FacetEngine) GetFacets() (map[string]*FacetGroup, error) {
	then := now()
	facetGroups := map[string]*FacetGroup{}
	for _, genericObject := range f.genericObjects {
		idPaths := f.facetPath.IDDotNotation
		if idPaths == "" {
			idPaths = "id"
		}
		id := getAtPathString(genericObject, strings.Split(idPaths, "."))
		if strings.TrimSpace(id) == "" {
			return nil, fmt.Errorf("found record with no id")
		}
		if f.initialized && !f.ids.Contains(id) {
			continue
		}
		f.allIds.Add(id)
		arrayPaths := strings.Split(f.facetPath.ArrayDotNotation, ".")
		namePaths := strings.Split(f.facetPath.NameFieldDotNotation, ".")
		nameMetaPaths := strings.Split(f.facetPath.NameMetaDotNotation, ".")
		valuePaths := strings.Split(f.facetPath.ValueMapDotNotation, ".")
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
				f.RecordLookup.Add(fmt.Sprintf("%s - %s", key, facetKey), &Record{
					ID:    id,
					Value: v,
				})
				facetGroup := facetGroups[key]
				if _, ok := facetGroup.Facets[facetKey]; !ok {
					facetGroup.Facets[facetKey] = &Facet{
						Name:   facetKey,
						Values: NewSet(),
					}
				}

				_, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, err
				}
				facetGroup.Facets[facetKey].Values.Add(v)
			}
		}
	}
	fmt.Printf("facet time: %d\n", now()-then)
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
