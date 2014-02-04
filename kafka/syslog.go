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
syslog harvester:

starting with file tail input:
* assume a file where new messages are appended and the file is rotated
* best effort to handle: deleted file, deleted lines at the top/middle of the file, symlinks and symlink changes???
* deduplicate messages by tracking line number and timestamp????
  XXX integrate this with kafka's "offset" ????
  kafka broker returns the "offset" with producer_response message



parse enough of the message to determine:
 - high precision time of the event
 - if it should also be sent through the low-latency path
 - if any metrics should be generated and sent directly to graphite



telling rsyslog to output useful fields:

  # http://www.rsyslog.com/what-is-the-difference-between-timereported-and-timegenerated/
  # can't rely on timereported. timegenerated is always the high precision time when syslog receives the message

  /etc/rsyslog.d/0-harvest.conf
  $template harvest,"%timegenerated:::date-rfc3339%,%syslogpriority-text%,%syslogfacility-text%,%programname%,%msg%\n"
  *.* /var/log/harvest;harvest

stage1:
 - tail fail and setup test cases

stage2:
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

follow_name behavior:
 - track the file by name.  handle inotify events for the file being moved, removed, renamed, written over, etc...





*/
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

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")

	// starting from logstash-forwarder/harvester.go

}
