
## Go Implementation details
The following go code will generate a facet group, with facets with facet values.

```golang
facetGroup, _ := CreateFacets(jsonData, &FacetPath{
  ArrayDotNotation:     "measurements",
  NameFieldDotNotation: "measurementName",
  NameMetaDotNotation:  "metrics.metricName",
  ValueMapDotNotation:  "metrics.measurements",
})
```

The return map of FacetGroups will have the following structure (json marshalled for viewing)

```json
{
  "FacetGroup": {
    "area (cube)": {
      "Name": "area (cube)",
      "Facets": {
        "side": {
          "Name": "side",
          "Values": [
            "10",
            "20"
          ]
        }
      }
    }
  }
}
```

You can then filter these results

```go
query := &Query{}
query.AddFilter("area (cube)", "side", Inclusive(8), Exclusive(12))
listOfIds := facetGroups.query(query)
```
