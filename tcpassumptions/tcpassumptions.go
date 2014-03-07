/*
package tcpassumptions attempts to establish reliable tcp client and server methods
that handle edge conditions????

use case:
 - client accumulates log events in a local file (length-prefix delimited + EOE marker)
 - client connects to server and dumps the whole file onto the network with sendfile(2)
   http://golang.org/src/pkg/net/sendfile_linux.go

use case:
 - client does gzip on send
 - server unzips, separates events and sends them on a channel for processing

reviewing tcp in detail just in case:

http://jan.newmarch.name/go/socket/chapter-socket.html

Once a client has established a TCP address for a service, it "dials" the service. If succesful, the dial returns a TCPConn for communication. The client and the server exchange messages on this. Typically a client writes a request to the server using the TCPConn, and reads a response from the TCPConn. This continues until either (or both) sides close the connection.


Implement connection limit?????
  track number of connections
  refuse connections if connection count is too high.


int listen(int sockfd, int backlog);
  defaults to /proc/sys/net/core/somaxconn which defaults to 128
  The backlog argument defines the maximum length to which the queue of pending connections for sockfd may grow.
  If a connection request arrives when the queue is full,
  the client may receive an error with an indication of ECONNREFUSED or,
  if the underlying protocol supports retransmission,
  the request may be ignored so that a later reattempt at connection succeeds.

can golang listen() pass a value for backlog????
  NO????? can't find it


tcplistener, err = ListenTCP(net, la)
  goroutine for each new connection
    read one message at a time and put the message on a sorting channel.
  block till enough work has been done????






http://synflood.at/tmp/golang-slides/mrmcd2012.html
*/
package main

import (
	"log"
	"net"
)

// Listen to tcp and divide the tcp into events
func handleConnection(c net.Conn) {
	// # make a buffer to read the event length

	// size the second buffer to read the whole event
	buf := make([]byte, 4096)

	for {
		n, err := c.Read(buf)
		if err != nil || n == 0 {
			c.Close()
			break
		}
		n, err = c.Write(buf[0:n])
		if err != nil {
			c.Close()
			break
		}
	}
	log.Printf("Connection from %v closed.", c.RemoteAddr())
}

func main() {

	// linux will queue up to 128 incoming connections
	// after that, clients will receive connection refused
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// handle one at a time and make everyone else wait???

		// make a buffered channel of connections and let a pool of goroutines read from them
		// block when the channel is full and let linux send connection refused
		// NOT go handleConnection(conn)
		some_chan <- conn
	}
}
