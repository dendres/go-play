package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os" // for File and friends
	"time"
)

/*

telling rsyslog to output useful fields:

  # http://www.rsyslog.com/what-is-the-difference-between-timereported-and-timegenerated/
  # can't rely on timereported. timegenerated is always the high precision time when syslog receives the message

  /etc/rsyslog.d/0-harvest.conf
  $template harvest,"%timegenerated:::date-rfc3339%,%syslogpriority-text%,%syslogfacility-text%,%programname%,%msg%\n"
  *.* /var/log/harvest;harvest

  by default, rsyslog guarantees that every log line will start with the timestamp.
  http://www.rsyslog.com/doc/rsconf1_escapecontrolcharactersonreceive.html

  with $MaxMessageSize 64k, got a message with 65556 characters... the remainder was truncated
     XXX sent to syslog from go. FYI: the bsd "logger" splits messages into 1024 character chunks

  every VALID line will have exactly the format specified in the template

  every VALID message will be 1024 characters or less.

  large messages are split by syslog, or by the logger application?????

  some effort should be put into restoring the escape sequences ????  probably not



*/

/*
kafka topic = t-<environment_name>
* a single topic name for all events
* prefix for the project to avoid conflicts with other projects. a single letter and dash should be sufficient.
* arbitrary string "environment" or "domain" is how you divide up servers by risk and customer exposure or management group etc...
*/

/*
kafka partition count = 15

each topic is partitioned into P partitions and replicated by factor N
 - partitions spread load across brokers
 - a single partition must not be bigger than the disk available on that broker
 - If you configure multiple data directories, partitions will be assigned round-robin to data directories. Each partition will be entirely in one of the data directories. If data is not well balanced among partitions this can lead to load imbalance between disks.
 - brokers do not enforce which message goes in which topic or partition
 - producer/consumer must agree on how to generate topic and partition for each message sent/received
   XXX or not they can both agree to not care and choose a random partition!!!!
 - once a topic and partition have been chosen, brokers can be asked which server is Leader for the given partition

https://cwiki.apache.org/confluence/display/KAFKA/FAQ#FAQ-HowdoIchoosethenumberofpartitionsforatopic?
 - "Clusters with up to 10k total partitions are quite workable"
 - more partitions mean smaller writes and more memory needed for VFS buffering
 - less partitions mean less kafka servers and more files in a given FS tree.
 - each partition has a small zookeeper cost.
 - more partitions mean more consumer checkpointing

if the client is configured to send/receive to/from multiple partitions, then it must keep multiple open tcp connections

3 kafka servers is ideal
5 kafka servers is ok
avoid more than 5 kafka servers
so the partition count should be a common multiple of both 3 and 5

assume separate 4TB hdd's
assume 1k messages under a 2/1 compression ratio
assume a target message rate of 10Km/s
how many partitions are required to store 30 days of messages in a single topic?
  bytes rate/s   m    h    d   days   KB     MB     GB     TB   disk size = min number of partitions
  500 * 10000 * 60 * 60 * 24 * 30 / 1024 / 1024 / 1024 / 1024 / 4         = 2.964

what is the message rate limit with 15 partitions?
  part TB    GB     MB     KB  bytes  msgs days    h    m    s = messages per second
  15 * 4 * 1024 * 1024 * 1024 * 1024 / 500 / 30 / 24 / 60 / 60 = 50Km/s

what is the disk write bandwidth used at 50Km/s?
   msgs bytes     KB     MB = MB per second
  50000 * 500 / 1024 / 1024 = 23MB/s XXX probably not near correct

*/

/*
starting with file tail input:
* assume a file where new messages are appended
* the file may be rotated
* the file may be a symlink pointing to another file that gets rotated
* follow_name behavior = track the file by name
* handle inotify events for the file being moved, removed, renamed, written over, etc...
* deduplicate messages by tracking line number, timestamp, and some kind of cheap checksum for the message???


deduplicate lines each time a new file is tailed!
journal messages:
 - message_id ??  increment, timestamp, source, checksum, etc...
 - read a file backlog from a configurable starting point (lines_back, duration_back)
 - once the journal is established, align the journal with the log file to find the correct starting point

buffer to disk. why??
- buffer to memory... cause disk swapping under pressure
- buffer to disk and VFS cache uses memory when it is available

XXXX disk buffer in single file or multiple files???
 - http://www.advancedlinuxprogramming.com/alp-folder/alp-apB-low-level-io.pdf




*/
func file_tail() {
}

type Harvester struct {
	// full Path to the file being tailed
	Path   string
	Offset int64
	Fields map[string]string

	file *os.File
}

// attempt to open h.Path. retry indefinitely until the file can be opened???
// seek to h.Offset
// return the new h.file
func (h *Harvester) open() *os.File {
	for {
		var err error
		h.file, err = os.Open(h.Path)
		if err != nil {
			log.Println("retry in 5 after error opening", h.Path, err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	// TODO(sissel): Only seek if the file is a file, not a pipe or socket.
	if h.Offset > 0 {
		h.file.Seek(h.Offset, os.SEEK_SET)
	} else if *from_beginning {
		h.file.Seek(0, os.SEEK_SET)
	} else {
		h.file.Seek(0, os.SEEK_END)
	}

	return h.file
}

func (h *Harvester) readline(reader *bufio.Reader, eof_timeout time.Duration) (*string, error) {
	var buffer bytes.Buffer
	start_time := time.Now()
	for {
		segment, is_partial, err := reader.ReadLine()

		if err != nil {
			if err == io.EOF {
				time.Sleep(1 * time.Second) // TODO(sissel): Implement backoff

				// Give up waiting for data after a certain amount of time.
				// If we time out, return the error (eof)
				if time.Since(start_time) > eof_timeout {
					return nil, err
				}
				continue
			} else {
				log.Println(err)
				return nil, err // TODO(sissel): don't do this?
			}
		}

		// TODO(sissel): if buffer exceeds a certain length, maybe report an error condition? chop it?
		buffer.Write(segment)

		if !is_partial {
			// If we got a full line, return the whole line.
			str := new(string)
			*str = buffer.String()
			return str, nil
		}
	} /* forever read chunks */

	return nil, nil
}

func (h *Harvester) Harvest(output chan *FileEvent) {
	if h.Offset > 0 {
		log.Printf("Starting harvester at position %d: %s\n", h.Offset, h.Path)
	} else {
		log.Printf("Starting harvester: %s\n", h.Path)
	}

	h.open()
	info, _ := h.file.Stat() // TODO(sissel): Check error
	defer h.file.Close()
	//info, _ := file.Stat()

	var line uint64 = 0 // Ask registrar about the line number

	// get current offset in file
	offset, _ := h.file.Seek(0, os.SEEK_CUR)

	log.Printf("Current file offset: %d\n", offset)

	// TODO(sissel): Make the buffer size tunable at start-time
	reader := bufio.NewReaderSize(h.file, 16<<10) // 16kb buffer by default

	var read_timeout = 10 * time.Second
	last_read_time := time.Now()
	for {
		text, err := h.readline(reader, read_timeout)

		if err != nil {
			if err == io.EOF {
				// timed out waiting for data, got eof.
				// Check to see if the file was truncated
				info, _ := h.file.Stat()
				if info.Size() < offset {
					log.Printf("File truncated, seeking to beginning: %s\n", h.Path)
					h.file.Seek(0, os.SEEK_SET)
					offset = 0
				} else if age := time.Since(last_read_time); age > (24 * time.Hour) {
					// if last_read_time was more than 24 hours ago, this file is probably
					// dead. Stop watching it.
					// TODO(sissel): Make this time configurable
					// This file is idle for more than 24 hours. Give up and stop harvesting.
					log.Printf("Stopping harvest of %s; last change was %d seconds ago\n", h.Path, age.Seconds())
					return
				}
				continue
			} else {
				log.Printf("Unexpected state reading from %s; error: %s\n", h.Path, err)
				return
			}
		}
		last_read_time = time.Now()

		line++
		event := &FileEvent{
			Source:   &h.Path,
			Offset:   offset,
			Line:     line,
			Text:     text,
			Fields:   &h.Fields,
			fileinfo: &info,
		}
		offset += int64(len(*event.Text)) + 1 // +1 because of the line terminator

		output <- event // ship the new event downstream
	} /* forever */
}

/*
tailer:

* runs on every server that produces log and metric events
* log file tail: process lines starting at the end of the file. go back max 24 hours. journal by timestamp
* add any missing fields (hostname, process/service etc.., timestamp, some kind of message id )
* remove any duplicate fields
* arbitrarily buffer a few messages and gzip
* kafka topic: "t-<environment_name>", partition: rand_int % partition_count=15
* send to kafka over optional encrypted channel
* choose a new topic and reconnect every M minutes


stage1:
 - just read from the end of the file and write the failing test cases

stage2:
 - then design to cover the cases
 - parse specific syslog format

stage3:
 - track line number and last few timestamps... to try for once and only once reading each line from the file

develop a strategy based on:
 - http://git.savannah.gnu.org/cgit/coreutils.git/tree/src/tail.c
 - https://github.com/elasticsearch/logstash-forwarder harvester.go
 - https://github.com/mozilla-services/heka/blob/dev/logstreamer/filehandling.go
   LogStream, LogStreamLocation
 - https://github.com/NarrativeScience/Log.io/blob/master/src/harvester.coffee
 - https://github.com/howeyc/fsnotify/blob/master/fsnotify_linux.go
 - https://github.com/jdamick/kafka
 - https://github.com/Shopify/sarama

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")

	// starting from logstash-forwarder/harvester.go

}
