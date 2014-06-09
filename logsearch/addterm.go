package main

import (
	"code.google.com/p/leveldb-go/leveldb/db"
	"code.google.com/p/leveldb-go/leveldb/table"
	"flag"
	"log"
	"os"
)

/*
add some test terms to the leveldb term database

term -> []byte is null character separated tokens

bootstrap assumptions:
* static time range
* static search term
* return json

*/

func main() {

	term := flag.String("term", "", "lowercase single-word key term")
	token := flag.String("token", "", "case sensitive key:value token that would appear in an event")
	flag.Parse()

	log.Println("term =", term, ", token =", token)

	dbname := "terms.db"
	dbfs := db.DefaultFileSystem

	// https://godoc.org/code.google.com/p/leveldb-go/leveldb/table
	// Tables are either opened for reading or created for writing but not both?????????
	// A reader can be used concurrently. Multiple goroutines can call Find concurrently

	// A writer writes key/value pairs in increasing key order, and cannot be used concurrently. A table cannot be read until the writer has finished.

	if _, err := os.Stat(dbname); os.IsNotExist(err) {
		log.Println("creating missing db file =", dbname)
		_, err := dbfs.Create(dbname)
		if err != nil {
			log.Fatalln("error creating db:", err)
		}
	}

	f, err := dbfs.Open(dbname)
	if err != nil {
		log.Fatalln(err)
	}

	w := table.NewWriter(f, nil)

	hello := []byte("hello")
	world := []byte("world of crap containing the whole world")

	w.Set(hello, world, nil)
	w.Close() //must call to write to file
}
