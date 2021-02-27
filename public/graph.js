// The code for this graph has been adapted from https://www.d3-graph-gallery.com/graph/connectedscatter_multi.html
// The code for the mouseover tooltip has been adapted from https://bl.ocks.org/d3noob/a22c42db65eb00d4e369

// set the dimensions and margins of the graph
var margin = {top: 20, right: 130, bottom: 30, left: 60},
    width = 800 - margin.left - margin.right,
    height = 700 - margin.top - margin.bottom;

// append the svg object to the body of the page
var svg = d3.select("#graph")
  .append("svg")
    .attr("width", "100%")
    .attr("height", height + margin.top + margin.bottom)
  .append("g")
    .attr("transform",
          "translate(" + margin.left + "," + margin.top + ")");

function formatTime (date) {
    return new Date(date).toLocaleDateString()
}

function formatNumber(n) {
    return n.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

// Define the div for the tooltip
var tooltip = d3.select("body").append("div")
    .attr("class", "tooltip")
    .style("opacity", 0);

//Read the data
d3.json("/stats", function (data) {
    // Order the data chronologically.
    data.sort((a, b) => Date.parse(a.date) - Date.parse(b.date))

    var dataReady = [
        {
            name: "First doses",
            values: data.map(function (d) {
                return {time: Date.parse(d.date), value: d.first_dose}
            }),
        },
        {
            name: "Second doses",
            values: data.map(function (d) {
                return {time: Date.parse(d.date), value: d.second_dose}
            }),
        },
    ]

    // Add X axis --> it is a date format
    var x = d3.scaleTime()
      .domain(d3.extent(dataReady[0].values, (d) => d.time))
      .range([ 0, width ]);
    svg.append("g")
      .attr("transform", "translate(0," + height + ")")
      .call(d3.axisBottom(x));

    // Add Y axis
    let maxY = dataReady[0].values[dataReady[0].values.length-1].value;
    var y = d3.scaleLinear()
      .domain([0, maxY])
      .range([ height, 0 ]);
    svg.append("g")
      .call(d3.axisLeft(y));

    // Add the lines
    var line = d3.line()
      .x(function(d) { return x(+d.time) })
      .y(function(d) { return y(+d.value) })
    svg.selectAll("lines")
      .data(dataReady)
      .enter()
      .append("path")
        .attr("d", function(d){ return line(d.values) } )
        .style("stroke-width", 2)
        .style("stroke", "black")
        .style("fill", "none")

    // Add the points
    svg
      // First we need to enter in a group
      .selectAll("dots")
      .data(dataReady)
      .enter()
        .append('g')
      // Second we need to enter in the 'values' part of this group
      .selectAll("points")
      .data(function(d){ return d.values })
      .enter()
      .append("circle")
        .attr("cx", function(d) { return x(d.time) } )
        .attr("cy", function(d) { return y(d.value) } )
        .attr("r", 5)
        .attr("stroke", "white")
        .on("mouseover", function(d) {
            tooltip.transition()
                .duration(50)
                .style("opacity", .95);
            tooltip.html(formatTime(d.time) + "<br/>"  + formatNumber(d.value))
                .style("left", (d3.event.pageX) + "px")
                .style("top", (d3.event.pageY - 28) + "px");
        })
        .on("mouseout", function(d) {
            tooltip.transition()
                .duration(1000)
                .style("opacity", 0);
        });

    // Add a legend at the end of each line
    svg
      .selectAll("labels")
      .data(dataReady)
      .enter()
        .append('g')
        .append("text")
          .datum(function(d) { return {name: d.name, value: d.values[d.values.length - 1]}; }) // keep only the last value of each time series
          .attr("transform", function(d) { return "translate(" + x(d.value.time) + "," + y(d.value.value) + ")"; }) // Put the text at the position of the last point
          .attr("x", 12) // shift the text a bit more right
          .text(function(d) { return d.name; })
          .style("font-size", 15)

})
