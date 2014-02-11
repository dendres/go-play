package main

import (
	"github.com/ActiveState/tail"
	"log"
	"time"
)

/*
tailer:

* runs on every server that produces log and metric events
* buffering happens ONLY in the tailed files and in memory.
* tail files by buffered inotify and record state on disk (fsync) after each block is sent on the network
* state/path/to/file. mtime = stat. 16 character, zero padded base10 offset in file + 64 characters of last line read from the file
* up to the entire content of the file may be processed if the file is rotated
* low latency between inotify and in-memory buffer should prevent loss of the last messages before rotate.
* add any missing fields (hostname, process/service etc.., timestamp, some kind of message id )
* remove any duplicate fields
* keep a long-standing compressed and encrypted messaging channel to multiple servers.
* sleep and retry on server disconnect or error in compression, encryption, or tcp.
* clean disconnect and reconnect periodically


starting with file tail input:
* lines are appended to the file
* the file may be rotated
* the file may be a symlink pointing to another file that gets rotated
* follow_name behavior = track the file by name
* deduplicate messages

network pattern:
* single long tcp connection to one server, or multiple stateless UDP connections to multiple servers?
* udp means max event size = (MTU - header_overhead) / compression_ratio = message size
* udp is a mild pain for firewalls... stateful TCP is easier.
* tcp max event size is limited by desired buffering characteristics (time and disk space)
* simultaneous message to multiple servers seems expensive?  is it???
  what about separate small and large buffer kafka topics?
    basically double the space to catch a very small number of missing messages.
* so far, buffer by min(time || disk space), compress, and send to a single kafka topic seems most efficient on network

consider this network pattern:
* long standing compressed, encrypted, tcp connections
* tail, parse, serialize, MTU sized in-memory buffer and send to randomly chosen 2 of 3 distributed receiving servers
* receiving process buckets by timestamp
* receiving process writes to local buffer for log.io server
* low latency receivers buffer and dedup messages
* SPOF process reads all buckets one bucket ago, and sends to "recent"
* SPOF process waits 48 hours, collects from all receivers and archives a bucket
   all state is in s3 and bucket and local files on the 5 servers
* SPOF process prunes remotes after 30 - 45 days


benefits of multiple send:
* message redundancy and distribution is decided at client
* all SPOF processes are stateless and can be very easily configured to heartbeat and fail over.


Error on the side of duplicate messages during file tailing race conditions!
but make sure that duplicate_messages / total_messages < 1% and hopefully a lot less.
kafka consumers can assume unique key of(time + checksum + hostname) for deduplication purposes


go channel pub/sub ????
* http://rogpeppe.wordpress.com/2009/12/01/concurrent-idioms-1-broadcasting-values-in-go-with-linked-channels/





buffer to disk. why??
 - buffer to memory... cause disk swapping under pressure
 - buffer to disk and VFS cache uses memory when it is available

disk buffer in single file or multiple files???
 - http://www.advancedlinuxprogramming.com/alp-folder/alp-apB-low-level-io.pdf
 - appending to a single file should be ok for now... get moving


case: start for the first time on a new server with files that have never been read
* offset and first64 are missing
* read the whole file

case: start after being off for 10min
* offset and first64 exist
* go to offset, compare first64
   match: process from next line to end of file
   no match: seek backward one line at a time till match is found, or beginning of file
   process as normal

case: file is rotated
* offset and first64 exist
* seek to offset and read first64... first64 will be null string or 0 or something... anyway, not a match
* work back to the beginning of the file and process as normal


each tail
  after each round of successful file buffering
    must save:
      time when last line was read
      timestamp from the line (if available)
      offset of the last line successfully buffered
      first 64 bytes of the line
    write to: offsets/full/file/path  ...in the same serialization as buffer, wire, etc...

if the offset is 0 or does not exist
  read the whole file

if offset > 0
  go to offset, read up to the first 64 bytes of the line
  if the line_start strings match
    continue reading from the next line to the end of the file
  if the don't match
    go back one line and try again
      continue till a match is found or the beginning of the file is reached
        (file was rotated or truncated and there's no way to get the data)
        then start reading from the next line to the end of the file


read line, parse, serialize, and append to buffer message/sta/mp
 - buffering configuration: S seconds, M messages
 - new buffer files will be created AT MOST every 1 second and will be larger than M
 - stamp is lp32.Encode(buffer_creation_time.Now().Unix(), 5) truncating higher order bits
 - stamp loops after 12.1 days
 - directory is messages/sta/mp
   a new "sta" is created every 17 minutes
   and a new "mp" is created at most every 1 second (for low M and high message rate

 - this avoids read/write contention

a separate goroutine processes the files:
* find the oldest buffer file
* read file, compress and send to kafka
* XXX determine kafka receive guarantee etc... decide if the any local message buffer is required ?????








*/

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")

	file := "/var/log/harvest"

}

// https://github.com/ActiveState/tail

// t, err := tail.TailFile(file, tail.Config{ReOpen: true, Follow: true})
// if err != nil {
// 	log.Fatal(err)
// }

// for line := range t.Lines {
// 	log.Println(line.Text)
// }

// log.Println("waiting 90 seconds")
// <-time.After(90 * time.Second)
