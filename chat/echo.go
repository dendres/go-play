// credit??  I think this was a go tour example?

package main

import (
	"io"
	"log"
	"net"
)

const lA = "localhost:4000"

func main() {
	log.Print("starting")
	listener, err := net.Listen("tcp", lA)
	if err != nil {
		log.Fatal(err)
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(connection, connection)
	}
}
