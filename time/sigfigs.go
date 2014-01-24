package main

import (
	"log"
)

/*
determine the precision and error of the local system clock.

determine how many bytes are required to store the system time with reasonable accuracy

http://stackoverflow.com/questions/16740014/computing-time-in-linux-granularity-and-precision


# BUG: leaving this at the following premature conclusion:
#  - ns for single server event comparison. see linux clock_gettime()
#  - us for LAN comparison if PTP is configured correctly and tested. over over LAN can do +-1us
#  - ms for internet: ntp over the public internet gets you at best +-10ms


# a logging system should allow a variable precision time to address multiple use cases



*/
func main() {
	log.Println("started ticker")
}
