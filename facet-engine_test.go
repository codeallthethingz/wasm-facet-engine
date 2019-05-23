package engine

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	facetGroups, err := CreateFacetGroups(readmeExample, &FacetPath{
		ArrayDotNotation:     "measurements",
		NameFieldDotNotation: "measurementName",
		NameMetaDotNotation:  "metrics.metricName",
		ValueMapDotNotation:  "metrics.measurements",
	})
	if err != nil {
		panic(err)
	}
	query := &Query{}
	query.AddFilter("area (cube)", "side", Inclusive(8), Exclusive(12))
	listOfIds, err := facetGroups.query(query)
	if err != nil {
		panic(err)
	}
	require.Equal(t, []string{"record 1"}, listOfIds)
}

func TestForIds(t *testing.T) {
	_, err := CreateFacetGroups("["+object8+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	require.Error(t, err)
	_, err = CreateFacetGroups("["+object8+"]", &FacetPath{
		IDDotNotation:        "name",
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	require.Nil(t, err)

}
func TestOptionalMetaName(t *testing.T) {
	// TODO
}

func TestQueryEdges(t *testing.T) {
	// TODO
	t.Fail()
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
	require.Equal(t, "{\"FacetGroup\":{\"area (cube)\":{\"Name\":\"area (cube)\",\"Facets\":{\"side\":{\"Name\":\"side\",\"Values\":[\"10\",\"20\"]}}}}}", string(data))
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
	require.Equal(t, 0, facetGroups.Len())
	facetGroups, err = CreateFacetGroups("["+object7+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "container.boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, facetGroups.Len())
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
	require.Equal(t, 0, facetGroups.Len())
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
	require.Equal(t, 0, facetGroups.Len())
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
	require.Equal(t, 0, facetGroups.Len())
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
	require.Equal(t, 1, facetGroups.Len())
	require.Equal(t, 3, len(facetGroups.Get("total-area (hex-cylinder)").Facets))
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
	require.Equal(t, 3, facetGroups.Len())
	require.Equal(t, 3, len(facetGroups.Get("total-area (hex-cylinder)").Facets))
	require.ElementsMatch(t, []string{"15", "16"}, facetGroups.Get("total-area (hex-cylinder)").Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"20"}, facetGroups.Get("total-area (hex-cylinder)").Facets["height"].Values.ToArray())
	require.Equal(t, []string{"1"}, facetGroups.Get("total-area (hex-cylinder)").Facets["weird"].Values.ToArray())
	require.Equal(t, 2, len(facetGroups.Get("head (hex-cylinder)").Facets))
	require.Equal(t, 3, len(facetGroups.Get("shaft (screwthread)").Facets))
}

func TestEdges(t *testing.T) {
	facetGroups := &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	}
	_, err := CreateFacetGroups("[{\"id\":\"1\"}]", facetGroups)
	require.Nil(t, err)
	_, err = CreateFacetGroups("[]", facetGroups)
	require.Nil(t, err)
	_, err = CreateFacetGroups("  ", facetGroups)
	require.Nil(t, err)
	_, err = CreateFacetGroups("NOTJSON", facetGroups)
	require.Error(t, err)
	_, err = CreateFacetGroups("["+object8+"]", facetGroups)
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
	require.Equal(t, 1, facetGroups.Len())
	require.Equal(t, 3, len(facetGroups.Get("total-area (hex-cylinder)").Facets))
	require.Equal(t, []string{"16"}, facetGroups.Get("total-area (hex-cylinder)").Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"1"}, facetGroups.Get("total-area (hex-cylinder)").Facets["weird"].Values.ToArray())
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
