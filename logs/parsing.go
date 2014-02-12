package main

import (
	"github.com/moovweb/rubex"
	"log"
	"time"
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

// type Event interface {} // some methods that all events must have

// line is the original plain string with date stamp etc... removed.
type Avent struct {
	time  time.Time
	shn   string // short hostname
	level string // log level or syslog priority
	app   string // overloaded = app_name || process_name || log_source || process_id
	line  string // the actual log line
}

// is there something like to_s in golang??? where if I override it, the print format changes?
func (a Avent) String() (out string) {
	out += a.time.String()
	// XXX finish making a readable way to print Avent
	return
}

// an event where the line has been replaced with an encoded byte slice
type Bvent struct {
	time  time.Time
	shn   string // short hostname
	level string // log level or syslog priority
	app   string // overloaded = app_name || process_name || log_source || process_id
	enc   []byte // the tokenized, encoded log line
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
	format := "2006-01-02T15:04:05.000000-07:00"

	test_string := `2014-01-30T00:30:01.246899+00:00,info,cron,CROND  (root) CMD (/usr/lib64/sa/sa1 1 1)`

	re := rubex.MustCompile(`^(\d+\-\d+\-\d+T\d+\:\d+\:\d+\.\d+\+\d+\:\d+),([^\,]+),([^\,]+),(.*)$`)
	matches := re.FindAllStringSubmatch(test_string, -1) // [][]string
	if matches != nil {
		match := matches[0]
		t, err := time.Parse(format, match[1])
		if err == nil {
			a := Avent{t, `fake-hostname`, match[2], match[3], match[4]}
			log.Println("found match =", a)
		} else {
			log.Println("time parsing error =", err)
		}
	}

	// FindAllStringSubmatch(s string, n int) [][]string

	log.Println("end main")
}

/*


    s := "2014-01-30T00:30:01.246899+00:00"
    time.Parse(format, s)
    fmt.Println(s)


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
