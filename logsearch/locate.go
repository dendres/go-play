package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/go-martini/martini"
	"io"
	"log"
	"net/http"
	"strings"
)

// assumes an html body containing a json array of key strings
// decodes the body and looks up each string
// returns map[string][]byte
func BoltMap(db_name string, bucket_name string, r io.Reader) (map[string][]byte, error) {
	m := make(map[string][]byte)
	keys := make([]string, 50)

	err := json.NewDecoder(r).Decode(&keys)
	if err != nil {
		return m, fmt.Errorf("error decoding json request: %v", err)
	}

	db, err := bolt.Open(db_name, 0666)
	if err != nil {
		return m, fmt.Errorf("error opening db: %v", err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket_name))

		for _, key := range keys {
			// values must be copied or they dissappear when the db is closed!
			v1 := bucket.Get([]byte(key))
			v2 := make([]byte, len(v1))
			copy(v2, v1)
			m[key] = v2
		}
		return nil
	})

	return m, nil
}

// Space Separated Values Converted to Json
// assumes each map value is really a string containing space separated values
// convert values from []byte to []string
// values may be nil!!
// json encode
// returns the html response code in the event of an error
func SSVJ(incoming map[string][]byte) (int, string) {
	out := make(map[string][]string)

	for k, v := range incoming {
		sv := string(v)
		out[k] = strings.Split(sv, " ")
	}

	js, err := json.Marshal(out)
	if err != nil {
		return 500, fmt.Sprintf("json.Marshal error: %v", err)
	}

	return 200, string(js)
}

// martini appears to run a pool of 6 handlers.
// probably need to serialize leveldb reads through channel to single goroutine
func main() {
	m := martini.Classic()
	m.Post("/terms", func(r *http.Request) (int, string) {

		log.Println("got body =", r.Body)

		m, err := BoltMap("terms.db", "terms", r.Body)
		if err != nil {
			log.Println(err)
			return 500, err.Error()
		}

		log.Println("calling ssvj")

		return SSVJ(m)
	})

	m.Post("/tokens", func(r *http.Request) (int, string) {

		m, err := BoltMap("terms.db", "tokens", r.Body)
		if err != nil {
			log.Println(err)
			return 500, err.Error()
		}

		log.Println("token map =", m)

		// get heka to logstream from /var/log/secure
		// and write normalized messages to file buffers
		// write service to read file buffers and write to term and token indexes

		// implement intersection!
		// then return array of hours

		return SSVJ(m)
	})

	log.Fatal(http.ListenAndServe(":8080", m))
}

/*
find all events where "Error", "Java", and "Exception" all appear anywhere in the same event over the last 7 years


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
