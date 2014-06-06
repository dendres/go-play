// Package ldbm implements aleveldb lookup via a request channel.
// leveldb keys and values are []byte
package ldbm

import (
	"code.google.com/p/leveldb-go/leveldb/db"
	"code.google.com/p/leveldb-go/leveldb/table"
	"fmt"
)

type Request struct {
	key      []byte
	vchan chan Response
}

type Response struct {
	value []byte
	err error
}


// Open the given leveldb file and return a request channel.
// Return error if leveldb returns more than "max" values.
func Open(dbfile) (requestChan chan Request, err error){
	dbfs := db.DefaultFileSystem
	f, err := dbfs.Open(dbfile)
	if err != nil {
		return requestChan, fmt.Errorf("error opening dbfile = %v. error = %v", dbfile, err)
	}
	reader := table.NewReader(f, nil)

	go func() {
		for request := range requestChan {
			value, err := reader.Get(request.key, nil)
			if err != nil {
				request.vchan <- Response{values, err}
			}
			request.vchan <- Response{values, nil}
		}
	}

	return requestChan
}

// Lookup takes a list of terms and returns a list of tokens matching those terms
func Lookup(terms []string) []string {

}

func main() {
	dbname := "hello.db"

	f0, err := dbfs.Create(dbname)
	if err != nil {
		log.Fatalln("error creating db", err)
	}

	w := table.NewWriter(f0, nil)

	hello := []byte("hello")
	world := []byte("the whole world")

	w.Set(hello, world, nil)
	w.Close() //must call to write to file

	f1, err := dbfs.Open(dbname)
	if err != nil {
		log.Fatalln("error getting record", err)
	}

	r := table.NewReader(f1, nil)

	out, err := r.Get([]byte("hello"), nil)
	if err != nil {
		log.Fatalln("error getting record", err)
	}

	log.Println(string(out))
}
