import facetEngine from "@realitypackagemanager/wasm-facet-engine"

function createRecords(recordCount) {
  let records = [];
  let metrics = createMetrics()
  for (let i = 0; i < recordCount; i++) {
    let record = {
      "id": "record " + (i + 1),
      "measurements": []
    }
    for (let k = 0; k < parseInt((Math.random() * 5)) + 1; k++) {
      record.measurements.push({
        "measurementName": "a" + (parseInt(Math.random() * 10)),
        "metrics": replaceX(metrics[parseInt(Math.random() * metrics.length)])
      })
    }
    records.push(record)
  }
  return records;
}

function replaceX(metric) {
  let newMetric = {
    metricName: metric.metricName,
    measurements: {}
  };
  Object.keys(metric.measurements).forEach(function (key) {
    newMetric.measurements[key] = parseInt(Math.random() * 100)
  });
  return newMetric
}
function createMetrics() {
  let metrics = []
  metrics.push({
    "metricName": "cube",
    "measurements": {
      "side": "X"
    }
  })
  metrics.push({
    "metricName": "screwthread",
    "measurements": {
      "height": "X",
      "diameter": "X",
      "pitch": "X"
    }
  })
  metrics.push({
    "metricName": "sphere",
    "measurements": {
      "diameter": "X"
    }
  })
  metrics.push({
    "metricName": "cuboid",
    "measurements": {
      "width": "X",
      "height": "X",
      "length": "X"
    }
  })
  metrics.push({
    "metricName": "cylinder",
    "measurements": {
      "diameter": "X",
      "height": "X"
    }
  })
  metrics.push({
    "metricName": "hexcube",
    "measurements": {
      "diameter": "X",
      "height": "X"
    }
  })
  metrics.push({
    "metricName": "pentcube",
    "measurements": {
      "diameter": "X",
      "height": "X"
    }
  })
  metrics.push({
    "metricName": "polygon",
    "measurements": {
      "diameter": "X",
      "height": "X",
      "sidecount": "X"
    }
  })
  metrics.push({
    "metricName": "oval-cylinder",
    "measurements": {
      "radius-major": "X",
      "radius-minor": "X",
      "height": "X"
    }
  })
  metrics.push({
    "metricName": "tricube",
    "measurements": {
      "width-hypotenuse": "X",
      "width-opposite": "X",
      "width-adjacent": "X",
      "height": "X"
    }
  })
  metrics.push({
    "metricName": "clip",
    "measurements": {
      "lip": "X",
      "lip-height": "X",
      "lip-width": "X",
      "height": "X",
      "width": "X",
      "length": "X"
    }
  })
  return metrics
}

let recordCount = 10;

let queriesExecuted = 0
let t
let timing = {}
let recordCountTiming = recordCount

let records = createRecords(recordCount)

facetEngine.load(function () {
  let config = {
    arrayDotNotation: "measurements",
    nameFieldDotNotation: "measurementName",
    nameMetaDotNotation: "metrics.metricName",
    valueMapDotNotation: "metrics.measurements",
  }
  console.log('init starting')
  t = new Date().getTime()
  facetEngine.initializeObjects(JSON.stringify(config), JSON.stringify(records), callbackFacets)
  timing['initialization'] = new Date().getTime() - t
  query()
})

function randomList(list) {
  return list[parseInt(Math.random() * list.length)];
}
function randomMap(map) {
  let keys = Object.keys(map)
  let key = randomList(keys)
  return map[key]
}
function randomKey(map) {
  let keys = Object.keys(map)
  return randomList(keys)
}
function query() {
  let record = randomList(records)
  let measurement = randomList(record.measurements)
  let measurementKey = randomKey(measurement.metrics.measurements)
  let groupName = measurement.measurementName + ' (' + measurement.metrics.metricName + ')'
  let facetName = measurementKey
  t = new Date().getTime()
  facetEngine.addFilter(groupName, facetName, true, 0, false, 90)
  facetEngine.query(callbackRecords, callbackFacets)
  queriesExecuted++
  if (queriesExecuted < 3) {
    query()
  } else {
    console.log(timing)
  }
}

function callbackFacets(facetGroupsJson) {
  timing['callback facets (' + recordCountTiming + ')'] = new Date().getTime() - t
}

function callbackRecords(idsJson) {
  let recordCountTiming = JSON.parse(idsJson).length
  console.log(JSON.parse(idsJson).length)
}