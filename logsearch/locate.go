package main

import (
	"fmt"
	"time"
)

/*

locate(downcased_term_list, time_range) returns term -> token -> day_stamp -> link

start interacting with the term index:
* downcase search term provided
* lookup search term
* return list of term info

bootstrap assumptions:
* static time range
* static search term
* return json

choice of embedded, single-file db for terms:
* lookup operations are disk operations! the whole thing is NOT read into memory!
* key -> value. no search or table-scan required
* keys and values are arbitrarily sized []byte
* some compression would be nice
* corruption detection and partial recovery would be nice

so far:
* leveldb-go
* https://github.com/syndtr/goleveldb
* jbarham/go-cdb
* tiedot
* sqlite





*/
func main() {
	format := "2006-01-02T15:04:05.000000-07:00"

	times := []string{
		"1970-01-01T00:00:00.000000+00:00",
		//		"2014-01-30T00:30:01.246899+00:00",
	}

	for m := 1; m < 13; m++ {
		for d := 0; d < 31; d++ {
			x := fmt.Sprintf("2014-%02d-%02dT00:00:00.000000+00:00", m, d)
			times = append(times, x)
		}
	}

	for i, time_string := range times {
		t, err := time.Parse(format, time_string)
		if err != nil {
			fmt.Println("error parsing timestamp", err)
		}
		time_int := uint64(t.Unix())
		fmt.Println(i, time_string, t, time_int)
	}
}
