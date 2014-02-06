package main

import (
	"github.com/ActiveState/tail"
	"log"
	"time"
)

/*
starting with file tail input:
* lines are appended to the file
* the file may be rotated
* the file may be a symlink pointing to another file that gets rotated
* follow_name behavior = track the file by name
* deduplicate messages

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


read line, parse, serialize, and append to messages/<10stamp>
 - based on the system time when the process is ready to write the file, not the time from the message

a separate goroutine processes the file
read 10stamp file
compress and send to kafka
if it's not the last file, delete








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
