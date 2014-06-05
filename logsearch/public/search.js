
window.onload=function(){

  var search_input = d3.select("#search_input");
  var show_terms = d3.select("#show_terms");
  var send_terms = d3.select("#send_terms");
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
  send_terms.on("click", function(datum,index){


    var post_string = JSON.stringify({"Terms": terms});
    console.log("sending: " + post_string );

    d3.xhr("/terms")
      .header("Content-Type", "application/json")
      .post(post_string, function(error, data){
        console.log("error = " + error);
        var out = JSON.parse(data.responseText);
        console.log(out);
      });

  });
};


//  var chCode = ('charCode' in event) ? event.charCode : event.keyCode;
//  alert ("The Unicode character code is: " + chCode);


