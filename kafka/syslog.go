package main

import (
	"log"
	"net"
)

/*
syslog harvester:

watch file... may be a fifo or a normal file????

serialize and bucket by timestamp
compress and send to kafka

allow configurable message durability for: disk io vs. memory footprint needs ???
message journal and deduplication???

crypto message authenticity???

map level and facility to kafka topics???

map messages to kafka partitions???
 - client's responsibility to determine which partition to send which message to??
 - default: hash(key) % numPartitions

rough disk based queue process for a single "queue":
* /queue/<source>/<interval>.txt
* wait "delay", then read file and send as a single kafka message
* delete file after ack

telling rsyslog to output useful fields:

  # http://www.rsyslog.com/what-is-the-difference-between-timereported-and-timegenerated/
  # can't rely on timereported. timegenerated is always the high precision time when syslog receives the message

  /etc/rsyslog.d/0-harvest.conf
  $template harvest,"%timegenerated:::date-rfc3339%,%syslogpriority-text%,%syslogfacility-text%,%programname%,%msg%\n"
  *.* /var/log/harvest;harvest

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")

	addr, err := net.ResolveUDPAddr("udp", ":5512")
	if err != nil {
		log.Println("error resolving udp addr", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("error establishing udp listener on", addr, ":", err)
	}

	// XXX using nc as client, got 2 messages of 2048 when sending very long string
	// determine what the correct buffer size is, or find a way to use a variable length buffer?
	// or not use a buffer at all???
	buf := make([]byte, 8192)

	for {
		rlen, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("error reading from udp:", err)
		}

		// addr, err := conn.ReadFromUDP(buf[0:])
		log.Println("from", remote, "length", rlen, ", got message:", string(buf[:rlen]))
	}

}
