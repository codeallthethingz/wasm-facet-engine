package engine

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForIds(t *testing.T) {
	// TODO
}
func TestOptionalMetaName(t *testing.T) {
	// TODO
}

func TestQuery(t *testing.T) {
	// TODO
}

func TestQueryEdges(t *testing.T) {
	// TODO
}

func TestUnmarshalSet(t *testing.T) {
	decoded := map[string]*FacetGroup{}
	err := json.Unmarshal([]byte("{\"area (cube)\":{\"Name\":\"area (cube)\",\"Facets\":{\"side\":{\"Name\":\"side\",\"Values\":\"esplode\"}}}}"), &decoded)
	require.Error(t, err)
}

func TestMarshalSet(t *testing.T) {
	facetGroups, err := CreateFacetGroups(readmeExample, &FacetPath{
		ArrayDotNotation:     "measurements",
		NameFieldDotNotation: "measurementName",
		NameMetaDotNotation:  "metrics.metricName",
		ValueMapDotNotation:  "metrics.measurements",
	})
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(facetGroups)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	require.Equal(t, strings.ReplaceAll("{\"area (cube)\":{\"Name\":\"area (cube)\",\"Facets\":{\"side\":{\"Name\":\"side\",\"Values\":[\"10\",\"20\"]}}}}", "\t", "  "), string(data))
	decoded := map[string]*FacetGroup{}
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		panic(err)
	}
	require.ElementsMatch(t, []string{"10", "20"}, decoded["area (cube)"].Facets["side"].Values.ToArray())
}

func TestBadPath(t *testing.T) {
	facetGroups, err := CreateFacetGroups("["+object7+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "subobject.subobject.name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
	facetGroups, err = CreateFacetGroups("["+object7+"]", &FacetPath{
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
	facetGroups, err := CreateFacetGroups("["+object6+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
}
func TestMissingName(t *testing.T) {
	facetGroups, err := CreateFacetGroups("["+object6+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
}
func TestNoValues(t *testing.T) {
	facetGroups, err := CreateFacetGroups("["+object5+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facetGroups))
}
func TestDots(t *testing.T) {

	facetGroups, err := CreateFacetGroups("["+object4+"]", &FacetPath{
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

	facetGroups, err := CreateFacetGroups("["+object1+","+object2+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 3, len(facetGroups))
	require.Equal(t, 3, len(facetGroups["total-area (hex-cylinder)"].Facets))
	require.ElementsMatch(t, []string{"15", "16"}, facetGroups["total-area (hex-cylinder)"].Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"20"}, facetGroups["total-area (hex-cylinder)"].Facets["height"].Values.ToArray())
	require.Equal(t, []string{"1"}, facetGroups["total-area (hex-cylinder)"].Facets["weird"].Values.ToArray())
	require.Equal(t, 2, len(facetGroups["head (hex-cylinder)"].Facets))
	require.Equal(t, 3, len(facetGroups["shaft (screwthread)"].Facets))
}

func TestEdges(t *testing.T) {
	facetGroups := &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	}
	_, err := CreateFacetGroups("[{}]", facetGroups)
	require.Nil(t, err)
	_, err = CreateFacetGroups("[]", facetGroups)
	require.Nil(t, err)
	_, err = CreateFacetGroups("  ", facetGroups)
	require.Nil(t, err)
	_, err = CreateFacetGroups("NOTJSON", facetGroups)
	require.Error(t, err)
}

func TestCapitalization(t *testing.T) {
	facetGroups, err := CreateFacetGroups("["+object2+","+object3+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 1, len(facetGroups))
	require.Equal(t, 3, len(facetGroups["total-area (hex-cylinder)"].Facets))
	require.Equal(t, []string{"16"}, facetGroups["total-area (hex-cylinder)"].Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"1"}, facetGroups["total-area (hex-cylinder)"].Facets["weird"].Values.ToArray())
}

var object1 = `{
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

var readmeExample = `
[
  {
    "name": "record 1",
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
    "name": "record 2",
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
