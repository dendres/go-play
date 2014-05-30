

// http://bl.ocks.org/mbostock/3883245
// http://bl.ocks.org/mbostock/3884955
// http://bl.ocks.org/benvandyke/8459843

var margin = {top: 20, right: 20, bottom: 30, left: 50},
width = 500 - margin.left - margin.right,
height = 200 - margin.top - margin.bottom;

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

var line = d3.svg.line()
  .x(function(d) { return x(d.X); })
  .y(function(d) { return y(d.Y); });

var svg = d3.select("body").append("svg")
  .attr("width", width + margin.left + margin.right)
  .attr("height", height + margin.top + margin.bottom)
  .append("g")
  .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

var legend = d3.select("body").append("ul")
  .attr("class", "legend");




d3.json("regression.json", function(error, data) {

  console.log("error =", error)
  console.log("data =", data)

  var graphs = data["Graphs"];
  var all_points = [];

  graphs.forEach(function(graph){
    var data_points = graph["DataPoints"];
    var regression_points = graph["RegressionPoints"];
    all_points = d3.merge([all_points, data_points, regression_points]);
  });

  x.domain(d3.extent(all_points, function(d){return d.X}));
  console.log("x.domain =", x.domain())

  y.domain(d3.extent(all_points, function(d){return d.Y}));
  console.log("y.domain =", y.domain())

  svg.append("g")
    .attr("class", "x axis")
    .attr("transform", "translate(0," + height + ")")
    .call(xAxis);

  svg.append("g")
    .attr("class", "y axis")
    .call(yAxis);

  var colors = d3.scale.category10();
  if (graphs.length > 10) {
    colors = d3.scale.category20();
  }

  // draw all the graphs
  for (i=0; i < graphs.length; ++i) {
    var color = colors(i);
    var graph = graphs[i];
    var data_points = graph["DataPoints"];
    var regression_points = graph["RegressionPoints"];
    var graph_name = graph["Name"];
    var r2 = graph["RSquared"];
    r2 = Math.round(r2*10000)/10000; // trim float

    svg.selectAll(".dot")
      .data(data_points)
      .enter().append("circle")
      .attr("r", 3)
      .attr("cx", function(d) { return x(d.X); })
      .attr("cy", function(d) { return y(d.Y); })
      .attr("fill", color);

    svg.append("path")
      .datum(regression_points)
      .attr("class", "line")
      .attr("stroke", color)
      .attr("d", line);

    legend.append("li")
      .text(graph_name +", r\u00B2 = "+ r2)
      .style("font-size", "22px")
      .style("list-style-type", "none")
      .style("color", color);
  }
});

