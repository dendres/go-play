package main

import (
	"code.google.com/p/leveldb-go/leveldb/db"
	"code.google.com/p/leveldb-go/leveldb/table"
	"log"
)

/*
add a term to the leveldb term database:

bootstrap assumptions:
* static time range
* static search term
* return json

*/

func main() {
	dbname := "hello.db"
	dbfs := db.DefaultFileSystem

	f0, err := dbfs.Create(dbname)
	if err != nil {
		log.Fatalln("error creating db", err)
	}

	w := table.NewWriter(f0, nil)

	hello := []byte("hello")
	world := []byte("world of crap containing the whole world")

	w.Set(hello, world, nil)
	w.Close() //must call to write to file

	f1, _ := dbfs.Open(dbname)
	r := table.NewReader(f1, nil)

	out, err := r.Get([]byte("hello"), nil)
	if err != nil {
		log.Fatalln("error getting record", err)
	}

	log.Println(string(out))
}
