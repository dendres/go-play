

// on each new character, list the terms
var show_terms = function(event) {
  var terms_string = d3.select("#search").property("value");
  var terms = terms_string.split(" ");
  console.log(terms);

  // populate the form with terms

};

var send_terms = function(event) {
  event.preventDefault();

  // on submit, send the terms to lookup handler

  console.log("would submit");
  return cancelDefaultAction(event);
};



//  var chCode = ('charCode' in event) ? event.charCode : event.keyCode;
//  alert ("The Unicode character code is: " + chCode);


