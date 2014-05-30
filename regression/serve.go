package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	d := "/home/done/go/src/github.com/dendres/go-play/regression/static"

	static := http.FileServer(http.Dir(d))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        static,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())

	//	panic(http.ListenAndServe(":8080", http.FileServer(http.Dir(d))))
}
