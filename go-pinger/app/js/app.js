'use strict';

var reload = setInterval(load, 5000)

function load() {
    fetch('/api/watchman')
        .then(res => res.json())
        .then(createCharts)
        .catch(console.error)
}

load()

function createCharts(data) {
    renderChart(data.filter(d => d.query_type === 'post'), '#chart1', 'count')
    renderChart(data.filter(d => d.query_type === 'post:featurized'), '#chart2', 'count')
    renderChart(data.filter(d => d.query_type === 'event'), '#chart3', 'count')
    renderChart(data.filter(d => d.query_type === 'jobset'), '#chart4', 'count')
    renderChart(data.filter(d => d.query_type === 'event'), '#chart5', 'resp_time')
}

function renderChart(data, chartSel, ofX) {
    d3.selectAll(chartSel + ' > *').remove();
    var chart = d3.timeseries()
        .addSerie(
        data,
        { x: 'created', y: ofX },
        { interpolate: 'linear', width: 3 }
        )
        .addHandlers([{ name: 'brush', cb: throttle(cancelReload, 1500) }])
        // .yscale.domain([0]) // show 0 on y axis
        .margin.left(65)
        .width(650)

    chart(chartSel)
}

// https://jsfiddle.net/jonathansampson/m7G64/
function throttle(callback, limit) {
    var wait = false;                  // Initially, we're not waiting
    return function () {               // We return a throttled function
        if (!wait) {                   // If we're not waiting
            callback.call();           // Execute users function
            wait = true;               // Prevent future invocations
            setTimeout(function () {   // After a period of time
                wait = false;          // And allow future invocations
            }, limit);
        }
    }
}

function cancelReload() {
    clearInterval(reload)
}
