package main

import (
	"github.com/moovweb/rubex"
	"log"
)

/*
event line parsing:

http://golang.org/pkg/regexp/
https://github.com/StefanSchroeder/Golang-Regex-Tutorial

https://github.com/moovweb/rubex and http://www.geocities.jp/kosako3/oniguruma/doc/RE.txt
* supposedly comparable syntax to, but faster than builtin
* yum install oniguruma-devel... for build
* claims to have named substring match, but I can't get it to work??!!!???


heka also includes https://github.com/losinggeneration/pego or something like it ????



how does capn protocol get data????  struct I believe????
see serialization.go



The FindAllStringSubmatch-function will, for each match, return an array with the entire match in the first field and the content of the groups in the remaining fields. The arrays for all the matches are then captured in a container array.

in the []byte, an entry will be "" if there is no match.

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")

	test_string := `2014-01-30T00:30:01.246899+00:00,info,cron,CROND  (root) CMD (/usr/lib64/sa/sa1 1 1)`

	re := rubex.MustCompile(`^(?<year>\d+)-(\d+)`)
	matches := re.FindAllStringSubmatch(test_string, -1)
	if matches == nil {
		log.Println("no match or error with match????")
	} else {
		log.Println("found matches =", matches)
	}

	// FindAllStringSubmatch(s string, n int) [][]string

	log.Println("end main")
}

/*

skip multiline for now:
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

*/
