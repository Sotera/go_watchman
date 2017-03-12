'use strict';

const reload = setInterval(load, 15000)

function load() {
    fetch('/api/watchman')
        .then(res => res.json())
        .then(createCharts)
        .catch(console.error)
}

load()

function createCharts(data) {
    // filtering is slow so use a worker
    [{
        queryType: 'post',
        chartSel: '#chart1',
        ofX: 'count'
    },
    {
        queryType: 'post:featurized',
        chartSel: '#chart2',
        ofX: 'count'
    },
    {
        queryType: 'event',
        chartSel: '#chart3',
        ofX: 'count'
    },
    {
        queryType: 'jobset',
        chartSel: '#chart4',
        ofX: 'count'
    },
    {
        queryType: 'event',
        chartSel: '#chart5',
        ofX: 'resp_time'
    }].forEach(obj => {
        // use 'script.999.js' format to match cache bust checker. can be any number.
        const worker = new Worker('/app/js/app-worker.999.js');
        worker.postMessage({ queryType: obj.queryType, data })
        worker.onmessage = event => {
            renderChart(event.data.data, obj.chartSel, obj.ofX)
        }
    })
}

function renderChart(data, chartSel, ofX) {
    d3.selectAll(chartSel + ' > *').remove();
    let chart = d3.timeseries()
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
    let wait = false;                  // Initially, we're not waiting
    return () => {               // We return a throttled function
        if (!wait) {                   // If we're not waiting
            callback.call();           // Execute users function
            wait = true;               // Prevent future invocations
            setTimeout(() => {   // After a period of time
                wait = false;          // And allow future invocations
            }, limit);
        }
    }
}

function cancelReload() {
    clearInterval(reload)
}
