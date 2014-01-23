package main

import (
	"fmt"
	"github.com/miekg/dns"
)

func main() {
	fmt.Println("Running...")
	m := new(dns.Msg)
	m.SetQuestion("google.com.", dns.TypeA)
	m.MsgHdr.RecursionDesired = true

	c := new(dns.Client)

	r, _, err := c.Exchange(m, "8.8.8.8:53")

	if err != nil {
		fmt.Println("got error result from exchange")
	}

	if r.Rcode != dns.RcodeSuccess {
		fmt.Println("invalid result from exchange")
	}

	fmt.Println(r)
}
