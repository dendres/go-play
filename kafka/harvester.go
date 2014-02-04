package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
harvester:
* runs on ever server that produces log and metric events
* input plugins accept various event serializations and normalize events enough to determine how to forward them
* add any missing fields (hostname, process/service etc.., timestamp, some kind of message id )
* remove any duplicate fields
* send to log.io server based on filter criteria
* buffer 30 seconds
* gzip
* kafka topic: "h-<environment_name>"
* kafka partition: based on timestamp... hash timestamp or mod timestamp or something????
* send to kafka


starting with file tail input:
* assume a file where new messages are appended
* best effort to handle: deleted file, deleted lines at the top/middle of the file
* the file may be rotated
* the file may be a symlink pointing to another file that gets rotated
* deduplicate messages by tracking line number and timestamp
* group messages into 30 second buckets


parse enough of the message to determine:
 - high precision time of the event
 - if it should also be sent through the low-latency path
 - if any metrics should be generated and sent directly to graphite

kafka topic:
 - h-<environment_name>

each topic is partitioned into P partitions and replicated by factor N
 - partitions spread load across brokers
 - a single partition must not be bigger than the disk available on that broker
 - brokers do not enforce which message goes in which topic or partition
 - producer/consumer must agree on how to generate topic and partition for each message sent/received
 - once a topic and partition have been chosen, brokers can be asked which server is Leader for the given partition

https://cwiki.apache.org/confluence/display/KAFKA/FAQ#FAQ-HowdoIchoosethenumberofpartitionsforatopic?
 - more partitions mean smaller writes and more memory needed for VFS buffering
 - less partitions mean less kafka servers and more files in a given FS tree.
 - each partition has a small zookeeper cost.
 - more partitions mean more consumer checkpointing

kafka partitioning:
 - timestamp: number of 30 second periods since unix epoch
 - message uniqueness: hostname + timestamp
 - small number of kafka brokers 3 or 5
 - not really a penalty for too many partitions????
 - make the partition_count a client setting... start with 100?

map hostname onto partitions:
 - hostname scheme: treat hostname as base36 integer
 - convert to base10 integer
 - % partition_count


so to send a message, the kafka client needs: environment and hostname ONLY!!!
these are converted to topic and partition






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
