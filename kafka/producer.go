package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
get events to their final format before they leave the server
split events for low-latency and high latency at the server
log events to a lossless long term data store for replay later
allow the events to be extracted by data_source and time_range for analysis
allow low-latency streaming of events for analysis in browser of incoming events
allow a dedicated high-performance recent-only index of the last 6 hours
allow the option to bypass kafka and send udp directly to the host(s) specified???
allow the option to have an unbuffered low latency kafka topic in addition to the long term storage topic
allow pruning by debug level?

message authenticity?????
* a bit of server identity
* how do I know this message is real?
* where did it come from?


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

kafka partitioning????
 - client's responsibility to determine which partition to send which message to??
 - default: hash(key) % numPartitions


log.io harvester: https://github.com/NarrativeScience/Log.io/blob/master/src/harvester.coffee
 - through kafka topic, or directly to a separate app????

metric creation:
 - fix the statsd bucketing errors and demonstrate flat heartbeat rate
 - send metrics through kafka like all other events for later replay, and/or??? directly to graphite

come up with a name for this app ????

deduplicate lines each time a new file is tailed!


TODO:
* read through "tail" and lumberjack to determine possible issues with file tailing
* test the go kafka producer and consumer
   https://github.com/jdamick/kafka
   https://github.com/Shopify/sarama
* investigate delivery guarantees, message identity requirements and message journaling
*  document how other projects do it

rough in-memory queue process for a single "source":
* parse and buffer incoming lines
* bucket messages into time intervals based on the time stamps and "interval" provided
* wait "delay", then compress and send the interval as a single kafka message
* delete bucket after kafka ack

rough disk based queue process for a single "queue":
* /queue/<source>/<interval>.txt
* wait "delay", then read file and send as a single kafka message
* delete file after ack





*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
