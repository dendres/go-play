
var svg = dimple.newSvg("#chartContainer", 590, 400);

d3.csv("data2.csv", function(error, data) {
  data.forEach(function(d) {
    d.xx = parseFloat(d.xaxis);
    d.yy = parseFloat(d.yaxis);
  });

  console.log(data);

  var myChart = new dimple.chart(svg, data);

  var x = myChart.addMeasureAxis("x", "xx");
  var y = myChart.addMeasureAxis("y", "yy");

  myChart.addSeries("thing", dimple.plot.line);
  myChart.addSeries("xx", dimple.plot.line);

  myChart.addLegend(200, 10, 360, 20, "right");
  myChart.draw();
});

