# Web Assembly Facet Engine

A piece of code designed to take a small set of json records and create facets based on sub parts of those records.
This is for a small number of records as it is intended to be used as a webassembly on a client browser.

The advantage of this approach is you can have very rich faceting without exploding your server.

The "small number" has yet to be determined, but performance metrics will be posted here when we know more.

**0.0.14 performance**

|Record Count|Initialization|Search
|--------:|----------:|--------:
| 1       | 2 ms      | 1 ms    
| 100     | 47 ms     | 2 ms    
| 1,000   | 528 ms    | 15 ms  
| 10,000  | 2,776 ms  | 42 ms
| 100,000 | 34,553 ms | 544 ms 

## Installation

*

Add the js bundle to your application

```bash
npm i @realitypackagemanager/wasm-facet-engine
```

Copy the `facet-engine.wasm` file from `node_modules/@realitypackagemanager/wasm-facet-engine` into the root of your web-application.

Ensure that this wasm file is served with the mime-type `application/wasm` or it will not work.

Import the facetEngine into your application

```node
import facetEngine from '@realitypackagemanager/wasm-facet-engine'
```

- `facetEngineLoad(callbackFunction)` - load the wasm file from your webserver. 
- `facetEngine.initializeObjects(stringifiedConfiguration, stringifiedObjectArray, callbackFacets)` - send in the records that you're going to work with and the configuration about which data elements are to be used as facets. Facets are sent back to the callback supplied
- `facetEngine.addFilter('facetGroupName', 'facetName', true, 7, false, 12)` - add a filter to the state.  The boolean parameters specify that the range is (true = inclusive) or (false = exclusive)
- `facetEngine.removeFilter(filterName)` - remove a filter by name
- `facetEngine.clearFilters()` - remove all filters
- `facetEngine.query(callbackRecords, callbackFacets)` - query the records for the current filters.  Results are sent to the supplied callback invocations `callbackFacets(stringifiedIdArray)`.  Facets are sent back to `callbackRecords(stringifiedFacets)`

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
facetEngine.load(function(){
  console.log('loaded!')
})
```

Initialize the engine with the list of objects to create facets for and the configuration of which data elements to extract for facets.

```javascript
let config = {
  arrayDotNotation:     "measurements",
  nameFieldDotNotation: "measurementName",
  nameMetaDotNotation:  "metrics.metricName",
  valueMapDotNotation:  "metrics.measurements"
}

facetEngine.initializeObjects(JSON.stringify(config), JSON.stringify(jsonData), function(facets){
  console.log('got facets', facets)
})
```

Add a filter and run it

```javascript
facetEngine.addFilter("area (cube)", "side", true, 8.0, false, 12.0)
facetEngine.query(callbackResults, callbackFacets)
function callbackResults(stringifiedIdArray) {
  listOfIds = JSON.parse(stringifiedIdArray)
  // Do something with the ids
}
function callbackFacets(stringifiedFacets) {
  facets = JSON.parse(stringifiedFacets)
  // Do something with the facets
}
```
