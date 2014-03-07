package main

import (
	"fmt"
)

/*

send to tcp: box(gzip(binary_event))
receive from tcp and unbox, unzip


each side generates: pk = crypto_box_keypair(&sk)


c = crypto_box(m,n,pk,sk)
where c is 16 bytes longer than m
m = crypto_box_open(c,n,pk,sk)


demonstrate nacl encrypt/decrypt channel over long standing TCP connection


http://cr.yp.to/highspeed/coolnacl-20120725.pdf
http://nacl.cr.yp.to/securing-communication.pdf
http://crypto.stackexchange.com/questions/6009/cryptographic-protocol-using-nacl

http://stackoverflow.com/questions/12741386/how-to-know-tcp-connection-is-closed-in-golang-net-package
http://jan.newmarch.name/go/socket/chapter-socket.html
http://commondatastorage.googleapis.com/io-2013/presentations/4053%20-%20Whispering%20Gophers.pdf


*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
