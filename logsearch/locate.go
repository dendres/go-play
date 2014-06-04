package main

import (
	"fmt"
	"time"
)

// find all events where "Error", "Java", and "Exception" all appear anywhere in the same event over the last 7 years

/*
Try 1: single data structure

each event has tokens:
  return events matching all tokens

token lookup: token -> []event_id
  return event_id's that appear in all sets with a single lookup

problems:
  single massive, constantly updating token -> []event_id table
  pruning by date requires table scan and rewrite
*/

/*
Try 2: map reduce

for each day:
  token -> []day_event_id
merge(all day_event_id sets)
7 year query requires 2555 index lookups and set operations
*/

/*
Try 3: token combinations return days

term:
 - case insensitive unicode string without whitespace or punctuation
 - the human input to the search process

combo:
 - sorted token combination
 - space separated
 - 1 to 4 tokens

token:
 - case sensitive and preserving
 - key_name:value_token

key_name:
 - strip non-printable characters
 - strip whitespace

value_token:
 - case sensitive and preserving
 - contains punctuation
 - no whitespace
 - only printable characters
 - extracted from fields by white space only

day_stamp:
 - seconds since epoch at the start of the day containing the event 00:00:00.0000000

indexes covering all 7 years: read-optimized b+ trees:
 - term -> token
 - token -> []combo
 - combo -> []day_stamp

indexes per day:
 - LSM or CDB? combo -> []event_id
 - token table
 - Binary pile of compressed events

search process:
 - autocomplete string to terms conversion in browser
 - submit terms and return tokens
 - select and submit 1-4 tokens, get back list of combos
 - submit combo, return list of days
 - submit day,combo, return event_id's
 - submit event_id, return event

daily prune process:
 - pick day to prune and get day_stamp
 - for each combo in the combo -> []event_id index
     find combo in combo -> []day_stamp index
       remove day_stamp
       if empty, delete combo
       else write new value for combo
       if empty, delete combo from term -> []combo index:
         extract terms from combo
         find all terms
         remove combo from any terms that contain the combo
         write new value for term
         if term is empty, delete term
 - delete per-day folder from S3 or object storage

daily index backup process:
 - boltdb snapshots or something?


*/

/*

start interacting with the term index:
* downcase search term provided
* lookup search term
* return list of term info

bootstrap assumptions:
* static time range
* static search term
* return json


first fake it with static data in the example below



choice of embedded, single-file db for terms:
* lookup operations are disk operations! the whole thing is NOT read into memory!
* key -> value. no search or table-scan required
* keys and values are arbitrarily sized []byte
* some compression would be nice
* corruption detection and partial recovery would be nice

Log Structured Merge-trees:
* leveldb-go
* https://github.com/syndtr/goleveldb

B+trees:
* sqlite
* boltdb: mmapped b+tree higher performance reads LSM, but slower writes

hash tables:
* jbarham/go-cdb
  fixed header
  24 bytes per record







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
