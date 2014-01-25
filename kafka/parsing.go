package main

import (
	"log"
)

/*
event line parsing:

https://github.com/moovweb/rubex
* supposedly faster than builtin http://golang.org/pkg/regexp/

multiline:
* start a multiline match when start_line_regex matches
  allow array of include_line_regex so that one or more regexes may look for lines that are part of the multiline
* stop when a line matches none of the include_line_regex
* stop after a timeout of N duration

multiline test cases:
* multiline json
* multiline xml
* java stack trace
* nodejs stack trace
* ruby stack trace

output message/event object format?
* ????


heka also includes https://github.com/losinggeneration/pego or something like it

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
