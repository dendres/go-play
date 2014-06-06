package main

import (
	"code.google.com/p/leveldb-go/leveldb/db"
	"code.google.com/p/leveldb-go/leveldb/table"
	"flag"
	"log"
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

	term := flag.String("term", nil, "lowercase single-word key term")
	token := flag.String("token", nil, "case sensitive key:value token that would appear in an event")
	flag.Parse()

	dbname := "terms.db"
	dbfs := db.DefaultFileSystem

	f1, err := dbfs.Open(dbname)
	// XXX if it's a not-exist error, then create the DB

	f0, err := dbfs.Create(dbname)
		if err != nil {
			log.Fatalln("error creating db", err)
		}

		w.Close() //must call to write to file


	}

	r := table.NewReader(f1, nil)

	w := table.NewWriter(f0, nil)

	hello := []byte("hello")

	// XXX replace world with null character separated tokens for testing

	world := []byte("world of crap containing the whole world")
	w.Set(hello, world, nil)





	out, err := r.Get([]byte("hello"), nil)
	if err != nil {
		log.Fatalln("error getting record", err)
	}

	log.Println(string(out))
}
