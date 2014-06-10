
// post json body and expect json body in response
// passes the response data to a callabck
var post = function(url, data, callback) {
  var post_string = JSON.stringify(data);
  console.log("posting: " + post_string );

  d3.xhr(url)
    .header("Content-Type", "application/json")
    .post(post_string, function(error, data){
      if (error) {
        console.warn(error);
        callback({});
      }
      var out = JSON.parse(data.responseText);
      callback(out);
    });
};

window.onload=function(){

  var search_div = d3.select("#search_div");
  var tokens_div = d3.select("#tokens_div");

  var search_input = d3.select("#search_input");
  var show_terms = d3.select("#show_terms");
  var send_terms = d3.select("#send_terms");
  var terms = [];

  // send tokens and handle the combos??????
  var sending_tokens = function(tokens) {
    post("/tokens", {"Tokens":tokens}, function(data){
      console.log(data);
      // XXXX setup the next round!!!!
    });
  };

  var sending_terms = function() {
    post("/terms", {"Terms": terms}, function(data){
      console.log(data);
      // clear out any old results
      d3.selectAll("#tokens_div *").remove();

      // consolidate tokens into a single list
      // XXX not merging for now. there will be duplicates
      tokens_list = [];
      for (var terms in data.Tokens) {
        var tokens = data.Tokens[terms];
        for (var i in tokens) {
          var token = tokens[i];
          tokens_list.push(token);
        }
      }
      console.log(tokens_list);

      // make a comma separated list/blob of tokens on the screen
      tokens_div.selectAll("span")
        .data(tokens_list)
        .enter()
        .append("span")
        .attr("class", "listed_token")
        .text(function(d,i){ return d; })
        .append("span")
        .text(", ");

      // make the token list multi-select
      d3.selectAll(".listed_token").on("click", function(datum, index){
        var li = d3.select(this);
        if (li.classed("selected_token")) {
          li.classed("selected_token", false);
        } else {
          if (d3.selectAll(".selected_token")[0].length < 4) {
            li.classed("selected_token", true);
          }
        }
      });

      // add Send Tokens button
      tokens_div.append("button")
        .attr("type", "button")
        .attr("id", "send_tokens")
        .text("Send Tokens")
        .on("click", function(datum, index){
          var tokens = [];
          d3.selectAll(".selected_token").each(function(d,i){
            tokens.push(d[0]);
          });

          sending_tokens(tokens);
        });
    });
  };

  // validate and display "terms" as they are typed into the input box
  search_input.on("keyup", function(datum,index) {
    terms = []; // validate terms
    search_input.property("value").split(" ").forEach(function(d,i){
      if (d.length > 3) {
        terms.push(d);
      }
    });

    if (terms.length > 0 ) {
      // show valid terms
      show_terms.text(terms.join(" "));
    }
  });

  // click send or hit enter to send terms
  send_terms.on("click", function(datum, index){
    sending_terms();
  });
  search_input[0][0].onkeydown = function (evt) {
    var keyCode = evt ? (evt.which ? evt.which : evt.keyCode) : event.keyCode;
    if (keyCode == 13) {
      sending_terms();
    }
  };


};


//  var chCode = ('charCode' in event) ? event.charCode : event.keyCode;
//  alert ("The Unicode character code is: " + chCode);


