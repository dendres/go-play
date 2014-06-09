package main

import (
	"flag"
	"github.com/boltdb/bolt"
	"log"
)

/*
add some test terms to the term database

term -> []byte is null character separated tokens

bootstrap assumptions:
* static time range
* static search term
* return json

*/

/*
https://godoc.org/github.com/boltdb/bolt

The DB is a collection of buckets and is represented by a single file on disk.
A bucket is a collection of unique keys that are associated with values.

read-only transactions
read-write transactions

*/
func main() {

	var term string
	flag.StringVar(&term, "term", "", "lowercase single-word key term")

	var token string
	flag.StringVar(&token, "token", "", "case sensitive key:value token that would appear in an event")

	flag.Parse()

	log.Println("adding: term =", term, ", token =", token)

	db, err := bolt.Open("terms.db", 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	bucket_name := []byte("terms")

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket_name)
		if err != nil {
			log.Fatalln(err)
		}

		return bucket.Put([]byte(term), []byte(token))
	})

	if err != nil {
		log.Fatalln("transaction failed:", err)
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket_name)
		return bucket.ForEach(func(k, v []byte) error {
			log.Printf("%s = %s.\n", k, v)
			return nil
		})
	})
}
