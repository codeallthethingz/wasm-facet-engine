# Web Assembly Facet Engine

A piece of code designed to take a small set of json records and create facets based on sub parts of those records.
This is for a small number of records as it is intended to be used as a webassembly on a client browser.

The advantage of this approach is you can have very rich faceting without exploding your server.

The "small number" has yet to be determined, but performance metrics will be posted here when we know more.

## Usage

Given the following array of two json objects held in a variable called `jsonData`:

```json
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
```

The following go code will generate a facet group, with facets with facet values.

```go
facetGroup, _ := CreateFacets(jsonData, &FacetPath{
  ArrayDotNotation:     "measurements",
  NameFieldDotNotation: "measurementName",
  NameMetaDotNotation:  "metrics.metricName",
  ValueMapDotNotation:  "metrics.measurements",
})
if err != nil {
  panic(err)
}
```

The return map of FacetGroups will have the following structure (json marshalled for viewing)

```json
{
  "area (cube)": {
    "Name": "area (cube)",
    "Facets": {
      "side": {
        "Name": "side",
        "Values": ["10", "20"]
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
