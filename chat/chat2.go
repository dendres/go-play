// http://synflood.at/tmp/golang-slides/mrmcd2012.html

package main

import (
	"log"
	"net"
)

func printMessages(msgchan <-chan string) {
	for msg := range msgchan {
		log.Printf("new message: %s", msg)
	}
}

func handleConnection(c net.Conn, msgchan chan<- string) {
	buf := make([]byte, 4096)

	for {
		n, err := c.Read(buf)
		if err != nil || n == 0 {
			c.Close()
			break
		}

		msgchan <- string(buf[0:n])

		n, err = c.Write(buf[0:n])
		if err != nil {
			c.Close()
			break
		}
	}
	log.Printf("Connection from %v closed.", c.RemoteAddr())
}

func main() {
	msgchan := make(chan string)
	go printMessages(msgchan)

	listener, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn, msgchan)
	}

}
