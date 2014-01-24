package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
log events to a lossless long term data store for replay later
allow the events to be extracted by data_source and time_range for analysis
allow low-latency streaming of events for analysis in browser of incoming events
allow a dedicated high-performance recent-only index of the last 6 hours

ensures all events including metric events get sent to kafka
 - allow durability / efficiency tradeoffs to be configurable per input
 - allow messages to be buffered by the message timestamp, or by system time

accept input from:
 - watch / tail a file and optionally manage it's rotation
 - syslog
 - a separate daemontools app to replace multilog
 - eventually, the rest of the logstash input formats

journal messages:
 - message_id ??  increment, timestamp, source, checksum, etc...
 - read a file backlog from a configurable starting point (lines_back, duration_back)
 - once the journal is established, align the journal with the log file to find the correct starting point

kafka output options:
 - some kind of header for figuring out which type of message this is???
 - buffer message by duration only. if the messages become too large, change the duration.


come up with a name for this app ????


deduplicate lines each time a new file is tailed!

same for lines received on udp

TODO:
* read through "tail" and lumberjack to determine possible issues with file tailing
* test the go kafka producer and consumer
   https://github.com/jdamick/kafka
   https://github.com/Shopify/sarama
* investigate delivery guarantees, message identity requirements and message journaling
*  document how other projects do it


*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
