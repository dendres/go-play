
window.onload=function(){

  var search_div = d3.select("#search_div");
  var tokens_div = d3.select("#tokens_div");

  var search_input = d3.select("#search_input");
  var show_terms = d3.select("#show_terms");
  var send_terms = d3.select("#send_terms");
  var tokens_found = d3.select("#tokens_found");
  var terms = [];

  // validate and display "terms" as they are typed into the input box
  search_input.on("keyup", function(datum,index) {
    terms = []; // validate terms
    search_input.property("value").split(" ").forEach(function(d,i){
      if (d.length > 4) {
        terms.push(d);
      }
    });

    if (terms.length > 0 ) {
      // show valid terms
      show_terms.text(terms.join(" "));
    }
  });

  // send the terms and display the matching tokens sent from the server
  send_terms.on("click", function(datum, index){

    var post_string = JSON.stringify({"Terms": terms});
    console.log("sending: " + post_string );

    d3.xhr("/terms")
      .header("Content-Type", "application/json")
      .post(post_string, function(error, data){
        console.log("error = " + error);
        var out = JSON.parse(data.responseText);
        console.log("server response = " + out);
        tokens_list = [];
        for (var k in out.Tokens) {
          var v = out.Tokens[k];
          tokens_list.push(v);
        }
        console.log("tokens_list = " + tokens_list);
        // now populate a multi-select list... ul/li where the elements get added to an input box
        tokens_found.selectAll("li")
          .data(tokens_list)
          .enter()
          .append("li")
          .attr("class", "listed_token")
          .text(function(d,i){ return d; });

        d3.selectAll(".listed_token").on("click", function(datum, index){
          console.log("clicked datum = " + datum);
          var li = d3.select(this);
          li.classed("selected_token", true);
        });

        // add send button
        tokens_div.append("button")
          .attr("type", "button")
          .attr("id", "send_tokens")
          .text("Send");

      });

  });



};


//  var chCode = ('charCode' in event) ? event.charCode : event.keyCode;
//  alert ("The Unicode character code is: " + chCode);


