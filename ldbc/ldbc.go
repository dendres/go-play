// Package ldbc implements a leveldb lookup via a request channel.
// leveldb keys and values are []byte
package ldbc

import (
	"code.google.com/p/leveldb-go/leveldb/db"
	"code.google.com/p/leveldb-go/leveldb/table"
	"fmt"
)

type Request struct {
	Key   []byte
	Vchan chan Response
}

type Response struct {
	Value []byte
	Err   error
}

// XXX need a way to retry reopening the DB on error etc...
// XXX should request and response be passed as pointers?????
// XXX need a way to keep a pool of reusable request/response objects to reduce garbage

// Open the given leveldb file and return a request channel.
// Return error if leveldb returns more than "max" values.
func Open(dbfile string) (requestChan chan Request, err error) {
	dbfs := db.DefaultFileSystem
	f, err := dbfs.Open(dbfile)
	if err != nil {
		return requestChan, fmt.Errorf("error opening dbfile = %v. error = %v", dbfile, err)
	}
	reader := table.NewReader(f, nil)

	go func() {
		for request := range requestChan {
			value, err := reader.Get(request.Key, nil)
			if err != nil {
				request.Vchan <- Response{value, err}
			}
			request.Vchan <- Response{value, nil}
		}
	}()

	return requestChan, nil
}
