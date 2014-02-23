package main

import (
	"fmt"
)

/*
rsyslog notes:
http://www.slideshare.net/rainergerhards1/rsyslog-vsjournal

systemd journal and the windows event log:
 - binary database
 - searchable (fast seek times)
 - ring buffer rollover
 - simple structured data inside the logs
 - collect metadata about sender: timestamp, uid, pid etc.. ?
 - unique id's for event classes?

"WinSyslog" "EventReporter" "snare for windows" "rsyslog agent for windows"

no single system is sufficient for all logging needs


http://tools.ietf.org/html/rfc5424
is there some standard for determining how much reliability is required for event delivery?

           Numerical         Severity
             Code

              0       Emergency: system is unusable
              1       Alert: action must be taken immediately
              2       Critical: critical conditions
              3       Error: error conditions
              4       Warning: warning conditions
              5       Notice: normal but significant condition
              6       Informational: informational messages
              7       Debug: debug-level messages

is it true that the more severe an event, the more delivery guarantee is required?
no. durability must be a separate field


disk_durability: 3 bits 0 - 7
  This is the number of different servers on which
  the event must be received
  and fsync() called on the file containing the event

verify_delivery: 0,1
0 = false: TCP ACK does not imply delivery

1 = true: TCP ACK implies delivery to the "disk_durability" specification

Urgency

durability: TCP will not ACK until:
0 = TCP ACK: disk local, disk datacenter, disk remote datacenter
1 = TCP ACK: disk local, disk datacenter, received remote datacenter
2 = TCP ACK: disk local, received datacenter, received remote datacenter
3 = TCP ACK: received local, received datacenter, received remote datacenter
4 = TCP ACK: received local, received datacenter
5 = TCP ACK: received local
6 = UDP: if received, it will be sent to more than one server via UDP or TCP
7 = UDP: if received, it will be sent to one other server via UDP




*/

func main() {
	fmt.Println("hi")
}
