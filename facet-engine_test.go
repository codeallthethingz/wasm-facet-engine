package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var defaultFacetPath = &FacetPath{
	ArrayDotNotation:     "bounds",
	NameFieldDotNotation: "name",
	NameMetaDotNotation:  "boundingType.name",
	ValueMapDotNotation:  "boundingType.measurements",
}
var exampleFacetPath = &FacetPath{
	ArrayDotNotation:     "measurements",
	NameFieldDotNotation: "measurementName",
	NameMetaDotNotation:  "metrics.metricName",
	ValueMapDotNotation:  "metrics.measurements",
}

func TestFilterForIds(t *testing.T) {
	facetEngine, _, err := NewFacetEngine(readmeExample, exampleFacetPath)
	if err != nil {
		panic(err)
	}
	facetEngine.ids.Add("record 1")
	facetGroups, err := facetEngine.GetFacets()
	require.Nil(t, err)
	require.ElementsMatch(t, []string{"10"}, facetGroups["area (cube)"].Facets["side"].Values.ToArray())

	facetEngine.ids = NewSet()
	facetEngine.ids.Add("bad record")
	facetGroups, err = facetEngine.GetFacets()
	require.Nil(t, err)
	require.Equal(t, 0, len(facetGroups))

}

func TestEmptyQuery(t *testing.T) {
	facetEngine, _, _ := NewFacetEngine("["+object9+"]", defaultFacetPath)
	_, _, err := facetEngine.Query(&Query{})
	require.Error(t, err)
}

func TestQueryEdges(t *testing.T) {
	_, _, err := NewFacetEngine("["+object9+"]", defaultFacetPath)
	require.Error(t, err)
}

func TestQueryInclusiveExclusive(t *testing.T) {
	testFilter(t, readmeExample, "area (cube)", "side", Inclusive(8), Inclusive(12), []string{"record 1"})
	testFilter(t, readmeExample, "area (cube)", "side", Inclusive(8), Inclusive(25), []string{"record 1", "record 2"})
	testFilter(t, readmeExample, "area (cube)", "side", Inclusive(10), Inclusive(20), []string{"record 1", "record 2"})
	testFilter(t, readmeExample, "area (cube)", "side", Exclusive(10), Exclusive(20), []string{})
	testFilter(t, readmeExample, "area (cube)", "side", Exclusive(10), Exclusive(10), []string{})
	testFilter(t, readmeExample, "area (cube)", "side", Exclusive(10), Inclusive(10), []string{})
	testFilter(t, readmeExample, "area (cube)", "side", Inclusive(10), Exclusive(10), []string{})
	testFilter(t, readmeExample, "area (cube)", "side", Inclusive(10), Inclusive(10), []string{"record 1"})
}
func TestQueryAnd(t *testing.T) {
	facetEngine, _, _ := NewFacetEngine("["+object1+","+object2+"]", defaultFacetPath)
	query := &Query{}
	query.AddFilter("total-area (hex-cylinder)", "height", Inclusive(0), Inclusive(25))
	listOfIds, _, _ := facetEngine.Query(query)
	require.ElementsMatch(t, []string{"1", "2"}, listOfIds)
	query.AddFilter("shaft (screwthread)", "pitch", Inclusive(1.5), Inclusive(1.5))
	listOfIds, _, _ = facetEngine.Query(query)
	require.ElementsMatch(t, []string{"1"}, listOfIds)
}

func testFilter(t *testing.T, example string, facetGroupName string, facetName string, min Range, max Range, expected []string) {
	facetEngine, _, err := NewFacetEngine(example, exampleFacetPath)
	if err != nil {
		panic(err)
	}
	listOfIds, _, err := facetEngine.Query((&Query{}).AddFilter(facetGroupName, facetName, min, max))
	require.Nil(t, err)
	require.ElementsMatch(t, expected, listOfIds)
}

func TestForIds(t *testing.T) {
	_, _, err := NewFacetEngine("["+object8+"]", defaultFacetPath)
	require.Error(t, err)
	_, _, err = NewFacetEngine("["+object8+"]", &FacetPath{
		IDDotNotation:        "name",
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	require.Nil(t, err)
}

func TestUnmarshalSet(t *testing.T) {
	decoded := map[string]*FacetGroup{}
	err := json.Unmarshal([]byte("{\"area (cube)\":{\"Name\":\"area (cube)\",\"Facets\":{\"side\":{\"Name\":\"side\",\"Values\":\"esplode\"}}}}"), &decoded)
	require.Error(t, err)
}

func TestMarshalSet(t *testing.T) {
	_, facetGroups, err := NewFacetEngine(readmeExample, exampleFacetPath)
	data, err := json.Marshal(facetGroups)
	if err != nil {
		panic(err)
	}
	require.Equal(t, "{\"area (cube)\":{\"Name\":\"area (cube)\",\"Facets\":{\"side\":{\"Name\":\"side\",\"Values\":[\"10\",\"20\"]}}}}", string(data))
	decoded := map[string]*FacetGroup{}
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		panic(err)
	}
	require.ElementsMatch(t, []string{"10", "20"}, decoded["area (cube)"].Facets["side"].Values.ToArray())
}

func TestBadPath(t *testing.T) {
	_, facetGroups, err := NewFacetEngine("["+object7+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "subobject.subobject.name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
	_, facetGroups, err = NewFacetEngine("["+object7+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "container.boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
}
func TestMissingMetaName(t *testing.T) {
	_, facetGroups, err := NewFacetEngine("["+object6+"]", defaultFacetPath)
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
}
func TestMissingName(t *testing.T) {
	_, facetGroups, err := NewFacetEngine("["+object6+"]", defaultFacetPath)
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
}
func TestNoValues(t *testing.T) {
	_, facetGroups, err := NewFacetEngine("["+object5+"]", defaultFacetPath)
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
}
func TestDots(t *testing.T) {
	_, facetGroups, err := NewFacetEngine("["+object4+"]", &FacetPath{
		ArrayDotNotation:     "container.bounds",
		NameFieldDotNotation: "container.name",
		NameMetaDotNotation:  "boundingType.container.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 1, len(facetGroups))
	require.Equal(t, 3, len(facetGroups["total-area (hex-cylinder)"].Facets))
}

func TestCreateFacets(t *testing.T) {

	_, facetGroups, err := NewFacetEngine("["+object1+","+object2+"]", defaultFacetPath)
	if err != nil {
		panic(err)
	}
	missing := facetGroups["missing value"]
	require.Nil(t, missing)

	require.Equal(t, 3, len(facetGroups))
	require.Equal(t, 3, len(facetGroups["total-area (hex-cylinder)"].Facets))
	require.ElementsMatch(t, []string{"15", "16"}, facetGroups["total-area (hex-cylinder)"].Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"20"}, facetGroups["total-area (hex-cylinder)"].Facets["height"].Values.ToArray())
	require.Equal(t, []string{"1"}, facetGroups["total-area (hex-cylinder)"].Facets["weird"].Values.ToArray())
	require.Equal(t, 2, len(facetGroups["head (hex-cylinder)"].Facets))
	require.Equal(t, 3, len(facetGroups["shaft (screwthread)"].Facets))
}

func TestEdges(t *testing.T) {
	facetPath := &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	}
	err := createFacetGroups("[{\"id\":\"1\"}]", facetPath)
	require.Nil(t, err)
	err = createFacetGroups("[]", facetPath)
	require.Nil(t, err)
	err = createFacetGroups("  ", facetPath)
	require.Nil(t, err)
	err = createFacetGroups("NOTJSON", facetPath)
	require.Error(t, err)
	err = createFacetGroups("["+object8+"]", facetPath)
	require.Error(t, err)
}

func createFacetGroups(data string, facetPath *FacetPath) error {
	_, _, err := NewFacetEngine(data, facetPath)
	return err
}

func TestCapitalization(t *testing.T) {
	_, facetGroups, err := NewFacetEngine("["+object2+","+object3+"]", defaultFacetPath)
	if err != nil {
		panic(err)
	}
	require.Equal(t, 1, len(facetGroups))
	require.Equal(t, 3, len(facetGroups["total-area (hex-cylinder)"].Facets))
	require.Equal(t, []string{"16"}, facetGroups["total-area (hex-cylinder)"].Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"1"}, facetGroups["total-area (hex-cylinder)"].Facets["weird"].Values.ToArray())
}

var object1 = `{
	"id": "1",
  "bounds": [
    {
      "name": "total-area",
      "boundingType": {
        "name": "hex-cylinder",
        "measurements": {
          "diameter": "15",
          "height": "20"
        }
      }
    },
    {
      "name": "head",
      "boundingType": {
        "name": "hex-cylinder",
        "measurements": {
          "diameter": "15",
          "height": "5"
        }
      }
    },
    {
      "name": "shaft",
      "boundingType": {
        "name": "screwthread",
        "measurements": {
          "diameter": "10",
          "height": "15",
          "pitch": "1.5"
        }
      }
    }
  ]
}`

var object2 = `{
	"id": "2",
  "bounds": [
    {
      "name": "total-area",
      "boundingType": {
        "name": "hex-cylinder",
        "measurements": {
          "diameter": "16",
          "height": "20",
          "weird": "1"
        }
      }
		}
	]
}`

var object3 = `{
	"id": "3",
  "bounds": [
    {
      "name": "total-AREA",
      "boundingType": {
        "name": "hex-cylinder",
        "measurements": {
          "diameter": "16",
          "height": "20",
          "Weird": "1"
        }
      }
		}
	]
}`
var object4 = `{
	"id": "4",
	"container" : {
		"bounds": [
			{
				"container" :  { "name": "total-AREA" }, 
				"boundingType": {
					"container": { "name": "hex-cylinder" },
					"measurements": {
						"diameter": "16",
						"height": "20",
						"Weird": "1"
					}
				}
			}
		]
	}
}`

var object5 = `{
	"id": "5",
  "bounds": [
    {
      "name": "total-AREA",
      "boundingType": {
        "name": "hex-cylinder",
        "measurements": {
        }
      }
		}
	]
}`
var object6 = `{
	"id": "6",
  "bounds": [
    {
      "boundingType": {
        "name": "hex-cylinder",
        "measurements": {
					"diameter": "16"
        }
      }
		}
	]
}`

var object7 = `{
	"id": "7",
  "bounds": [
    {
      "name": "total-AREA",
      "boundingType": {
				"name": "hex-cylinder",
        "measurements": {
					"diameter": "16"
        }
      }
		}
	]
}`

var object8 = `{
	"name": "8",
  "bounds": [
    {
      "name": "total-AREA",
      "boundingType": {
				"name": "hex-cylinder",
        "measurements": {
					"diameter": "16"
        }
      }
		}
	]
}`
var object9 = `{
	"id": "9",
  "bounds": [
    {
      "name": "total-AREA",
      "boundingType": {
				"name": "hex-cylinder",
        "measurements": {
					"height": "16h"
        }
      }
		}
	]
}`

var readmeExample = `
[
  {
    "id": "record 1",
    "measurements": [
      {
        "measurementName": "area",
        "metrics": {
          "metricName": "cube",
          "measurements": {
            "side": "10"
          }
        }
      }
    ]
  },
  {
    "id": "record 2",
    "measurements": [
      {
        "measurementName": "area",
        "metrics": {
          "metricName": "cube",
          "measurements": {
            "side": "20"
          }
        }
      }
    ]
  }
]
`
