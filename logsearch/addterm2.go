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

	var bucket string
	var key string
	var value string

	flag.StringVar(&bucket, "bucket", "", "bucket to work on")
	flag.StringVar(&key, "key", "", "key to set")
	flag.StringVar(&value, "value", "", "value to set")

	flag.Parse()

	if bucket == "" {
		log.Fatalln("missing -bucket argument")
	}

	if key == "" {
		log.Fatalln("missing -key argument")
	}

	if value == "" {
		log.Fatalln("missing -value argument")
	}

	log.Println("adding to bucket =", bucket, ": key =", key, ", value =", value)

	db, err := bolt.Open("terms.db", 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	bucket_name := []byte(bucket)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket_name)
		if err != nil {
			log.Fatalln(err)
		}

		return bucket.Put([]byte(key), []byte(value))
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
