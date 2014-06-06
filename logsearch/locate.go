package main

import (
	//	"encoding/json"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"io"
	"log"
	"net/http"
)

// find all events where "Error", "Java", and "Exception" all appear anywhere in the same event over the last 7 years

/*
Try 3: token combinations return days

term:
 - case insensitive unicode string without whitespace or punctuation
 - the human input to the search process

token:
 - case sensitive and preserving
 - key_name:value_token

combo:
 - sorted token combination
 - space separated
 - 1 to 4 tokens

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
 - term  -> []token
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

/*
first fake it with static data in the example below

* web form makes terms

pick web framework:
* martini
*


*/

type Terms struct {
	Terms []string
}

// Decode takes the http Request body and returns a slice of string terms.
func (t *Terms) Decode(body io.Reader) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&t)
	if err != nil {
		return fmt.Errorf("error decoding json request: %v", err)
	}

	return nil
}

type Response struct {
	Name   string
	Tokens []string
}

func main() {

	// go fireup the terms db handler and get a channel to send requests to
	// termsChan := terms.Open(dbfilename)

	m := martini.Classic()
	m.Post("/terms", func(r *http.Request) (int, string) {
		t := Terms{}

		err := t.Decode(r.Body)
		if err != nil {
			log.Println(err)
			return 500, err.Error()
		}

		log.Println(t.Terms)

		// tokenRequest := &terms.Request{[]int{3, 4, 5}, sum, make(chan int)}
		// termsChan <- tokenRequest
		// tokens := <- tokenRequest.tokensChan

		data := Response{
			"hello",
			[]string{"Token1", "FieldName:Token2", "id:555"},
		}

		b, err := json.Marshal(data)
		if err != nil {
			log.Println("json.Marshal error:", err)
		}
		return 200, string(b)
	})

	log.Fatal(http.ListenAndServe(":8080", m))
}
