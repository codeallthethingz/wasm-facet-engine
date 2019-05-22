package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateFacets(t *testing.T) {

	facets, err := CreateFacets("["+facetString1+","+facetString2+"]", &FacetPath{
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
	require.Equal(t, []string{"15", "16"}, facets["total-area (hex-cylinder)"].Facets["diameter"].Values.ToArray())
	require.Equal(t, []string{"20"}, facets["total-area (hex-cylinder)"].Facets["height"].Values.ToArray())
	require.Equal(t, []string{"1"}, facets["total-area (hex-cylinder)"].Facets["weird"].Values.ToArray())
	require.Equal(t, 2, len(facets["head (hex-cylinder)"].Facets))
	require.Equal(t, 3, len(facets["shaft (screwthread)"].Facets))
}

func TestNoValues(t *testing.T)       {}
func TestCapitalization(t *testing.T) {}

var facetString1 = `{
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

var facetString2 = `{
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
