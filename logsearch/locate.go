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

// Tokens returns the tokens for the given terms as a map
func (t *Terms) Tokens() (map[string][]string, error) {

	db, err := bolt.Open("terms.db", 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	tokens := make(map[string][]string)
	bucket_name := []byte("terms")

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket_name)

		for _, term := range t.Terms {
			token_bytes := bucket.Get([]byte(term))
			token_string := string(token_bytes)
			token_list := strings.Split(token_string, " ")
			log.Println("term =", term, "gave tokens =", token_list)
			tokens[term] = token_list
		}
		return nil
	})

	return tokens, nil
}

type Response struct {
	Name   string
	Tokens map[string][]string
}

// martini appears to run a pool of 6 handlers.
// probably need to serialize leveldb reads through channel to single goroutine
func main() {
	m := martini.Classic()
	m.Post("/terms", func(r *http.Request) (int, string) {
		t := Terms{}

		err := t.Decode(r.Body)
		if err != nil {
			log.Println(err)
			return 500, err.Error()
		}

		tokens, err := t.Tokens()
		if err != nil {
			log.Println("error looking up tokens in db:", err)
			return 500, err.Error()
		}

		data := Response{
			"tokens",
			tokens,
		}

		b, err := json.Marshal(data)
		if err != nil {
			log.Println("json.Marshal error:", err)
		}
		return 200, string(b)
	})

	log.Fatal(http.ListenAndServe(":8080", m))
}
