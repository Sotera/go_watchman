// def: a web worker to transform data

// importScripts('https://d3js.org/d3.v4.js');

onmessage = function (event) {
    let {data, queryType} = event.data;

    data = data.filter(d => d.query_type === queryType)

    postMessage({ data });
};
