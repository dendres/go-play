package main

import (
	"fmt"
)

/*
Log Event Data Scale Estimation:

event_rate_target:
* average events from all servers in one second?
* currently 6M events/hour/75servers ~ 22 events/second/server
* target rate is 500 events/second/server * 1K servers = 500K events/second
* 500K events/second = 500 events/ms = 0.5 events/us
* 500K * 512 byte events/second = 244MB/second = 20TB/day

bucket_event_count and storage size:
* current 2K events / second * 1024 seconds ~ 2M events/bucket
* target 500K events / second * 1024 seconds ~ 512M events/bucket
* and around 244MB/second * 1024 seconds ~ 244GB/bucket means 2 buckets per ec2 disk
* should be possible to handle the indexing for 512M event id's per bucket... right???
* assuming massive compression is possible...

event_key:
* how many bits of "point" are needed inside directory?
  32^2 = 1024 seconds = 2^10 = the lowest 10 bits
  fractional time = number of 10^-8 second (10ns) intervals since the second began.
  10^8, 10ns intervals = 1 second
  1 second or 10^8 10ns intervals fit in 27 bits
  so... 10 bits for seconds + 27 bits for fractional seconds = 37 bits for "point" inside bucket

* probability of collision using time stamp alone?
  this is less events than the granularity of Event.point = 10ns intervals
  calculate the probability of 2 messages at the same time????
    NTP cloud time accuracy ~ 10^-3 or 1ms
     every 1 ms can contain 10^5 * 10^-8 intervals or ~ 17 bits worth
    given our 37 bit "bucket point":
      20 bits can be considered collision free, but 17 bits must be considered like a random hash

* so how much extra entropy is required to avoid message collision in a "bucket"??
  total number of messages in 1ms = 5*10^5 events/second * 10^-3 seconds/ms = 500 events/ms
  bits of entropy required to avoid collision in 500 events:
    ((500)^3)/2 fits in ~ 26 bits
  we've got 17 bits already and need a total of 26. find 9 bits somewhere or risk it?
  XXXX I'm considering dropping the granularity of the timestamp in favor of checksum
    how many bits of timestamp in bucket for ms granularity?
      1 second or 10^3 ms fits in 2^10 or 10 bits
      so... 10bits for seconds + 10 bits for fractional seconds + 26 bits entropy = 46 bits required
      how to divide up the 6 bytes?  3 time + 3 checksum should be ok

Summary: 3 bytes of time + 3 bytes of checksum is the minimum required to avoid collision

*/

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
