# Web Assembly Facet Engine

A piece of code designed to take a small set of json records and create facets based on sub parts of those records.
This is for a small number of records as it is intended to be used as a webassembly on a client browser.

The advantage of this approach is you can have very rich faceting without exploding your server.

The "small number" has yet to be determined, but performance metrics will be posted here when we know more.

## Installation

Add the js bundle to your application

```bash
npm i @codeallthethingz/facet-engine
```

Download the `facet-engine.wasm` file from https://github.com/codeallthethingz/wasm-facet-engine/releases

Ensure that this wasm file is served with the mime-type `application/wasm` or it will not work.

This will add the following functions into the global scope
- `facetEngineLoad(facetEngineWasmLocation)` - load the wasm file into memory. Default: `facet-engine.wasm`
- `facetEngineInitializeObjects(stringifiedConfiguration, stringifiedObjectArray)` - send in the records that you're going to work with and the configuration about which data elements are to be used as facets. Facets are sent back to `facetEngineCallbackFacets(stringifiedFacets)`
- `facetEngineAddFilter('filterGroupName', 'metricName', true, 7, false, 12)` - add a filter to the state.  The boolean parameters specify that the range is (true = inclusive) or (false = exclusive)
- `facetEngineRemoveFilter(filterName)` - remove a filter by name
- `facetEngineClearFilters()` - remove all filters
- `facetEngineQuery()` - return a json array of id's that match all the set filters.  Results are sent to a callback invocation of a function you should add called `facetEngineCallbackResults(stringifiedIdArray)`.  Facets are sent back to `facetEngineCallbackFacets(stringifiedFacets)`

## Usage

Given the following array of two json objects held in a variable called `jsonData`:

```javascript
let jsonData = [
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

Load the facet engine

```javascript
facetEngineLoad()
```

Initialize the engine with the list of objects to create facets for and the configuration of which data elements to extract for facets.

```javascript
let config = {
  arrayDotNotation:     "measurements",
  nameFieldDotNotation: "measurementName",
  nameMetaDotNotation:  "metrics.metricName",
  valueMapDotNotation:  "metrics.measurements",
}
// this will call back to facetEngineCallbackFacets(stringifiedFacets)
facetEngineInitializeObjects(JSON.stringify(config), JSON.stringify(jsonData)) 
```

Add a filter and run it

```javascript
query.AddFilter()
facetEngineAddFilter("area (cube)", "side", true, 8.0, false, 12.0)
facetEngineQuery()
function facetEngineCallbackResults(stringifiedIdArray) {
  listOfIds = JSON.parse(stringifiedIdArray)
  // Do something with the ids
}
function facetEngineCallbackFacets(stringifiedFacets) {
  facets = JSON.parse(stringifiedFacets)
  // Do something with the facets
}
```
