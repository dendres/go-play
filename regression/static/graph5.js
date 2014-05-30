

// http://bl.ocks.org/mbostock/3883245
// http://bl.ocks.org/mbostock/3884955
// http://bl.ocks.org/benvandyke/8459843

var margin = {top: 20, right: 20, bottom: 30, left: 50},
width = 960 - margin.left - margin.right,
height = 500 - margin.top - margin.bottom;

var x = d3.scale.linear()
  .range([0, width]);

var y = d3.scale.linear()
  .range([height, 0]);

var xAxis = d3.svg.axis()
  .scale(x)
  .orient("bottom");

var yAxis = d3.svg.axis()
  .scale(y)
  .orient("left");

// d.XXXXX
var line = d3.svg.line()
  .x(function(d) { return x(d.X); })
  .y(function(d) { return y(d.Y); });

var svg = d3.select("body").append("svg")
  .attr("width", width + margin.left + margin.right)
  .attr("height", height + margin.top + margin.bottom)
  .append("g")
  .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

// regression.json has array of array of points
// each set of points must be plotted... all on the same graph
d3.json("regression.json", function(error, data) {

  console.log("got some data", data)

  // XXX have to set x and y domain over all data sets?????

  x.domain(d3.extent(data, function(d) { return d.X; }));
  y.domain(d3.extent(data, function(d) { return d.Y; }));

  svg.append("g")
    .attr("class", "x axis")
    .attr("transform", "translate(0," + height + ")")
    .call(xAxis);

  svg.append("g")
    .attr("class", "y axis")
    .call(yAxis);
    // .append("text")
    // .attr("transform", "rotate(-90)")
    // .attr("y", 6)
    // .attr("dy", ".71em")
    // .style("text-anchor", "end")
    // .text("Price ($)");

  svg.append("path")
    .datum(data)
    .attr("class", "line")
    .attr("d", line);
});

