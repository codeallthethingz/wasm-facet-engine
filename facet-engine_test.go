package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBadPath(t *testing.T) {
	facets, err := CreateFacets("["+object7+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "subobject.subobject.name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facets))
	facets, err = CreateFacets("["+object7+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "container.boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facets))
}
func TestMissingMetaName(t *testing.T) {
	facets, err := CreateFacets("["+object6+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facets))
}
func TestMissingName(t *testing.T) {
	facets, err := CreateFacets("["+object6+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facets))
}
func TestNoValues(t *testing.T) {
	facets, err := CreateFacets("["+object5+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, len(facets))
}
func TestDots(t *testing.T) {

	facets, err := CreateFacets("["+object4+"]", &FacetPath{
		ArrayDotNotation:     "container.bounds",
		NameFieldDotNotation: "container.name",
		NameMetaDotNotation:  "boundingType.container.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 1, len(facets))
	require.Equal(t, 3, len(facets["total-area (hex-cylinder)"].Facets))
}

func TestCreateFacets(t *testing.T) {

	facets, err := CreateFacets("["+object1+","+object2+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 3, len(facets))
	require.Equal(t, 3, len(facets["total-area (hex-cylinder)"].Facets))
	require.ElementsMatch(t, []string{"15", "16"}, facets["total-area (hex-cylinder)"].Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"20"}, facets["total-area (hex-cylinder)"].Facets["height"].Values.ToArray())
	require.Equal(t, []string{"1"}, facets["total-area (hex-cylinder)"].Facets["weird"].Values.ToArray())
	require.Equal(t, 2, len(facets["head (hex-cylinder)"].Facets))
	require.Equal(t, 3, len(facets["shaft (screwthread)"].Facets))
}

func TestEdges(t *testing.T) {
	facetPath := &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	}
	_, err := CreateFacets("[{}]", facetPath)
	require.Nil(t, err)
	_, err = CreateFacets("[]", facetPath)
	require.Nil(t, err)
	_, err = CreateFacets("  ", facetPath)
	require.Nil(t, err)
	_, err = CreateFacets("NOTJSON", facetPath)
	require.Error(t, err)
}

func TestCapitalization(t *testing.T) {
	facets, err := CreateFacets("["+object2+","+object3+"]", &FacetPath{
		ArrayDotNotation:     "bounds",
		NameFieldDotNotation: "name",
		NameMetaDotNotation:  "boundingType.name",
		ValueMapDotNotation:  "boundingType.measurements",
	})
	if err != nil {
		panic(err)
	}
	require.Equal(t, 1, len(facets))
	require.Equal(t, 3, len(facets["total-area (hex-cylinder)"].Facets))
	require.Equal(t, []string{"16"}, facets["total-area (hex-cylinder)"].Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"1"}, facets["total-area (hex-cylinder)"].Facets["weird"].Values.ToArray())
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
