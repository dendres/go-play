package hose

import (
	"fmt"
)

/*
Start the event sorting goroutine.
Singleton.
takes incoming events on a channel
keeps a trie of channels to file writing goroutines

watches a split channel that any writer can send to.
pauses event processing and updates the trie when a split is received


handle splits????
* writer has event_channel and split_channel
* writer sends True on split_channel when it's ready to split

https://github.com/timtadh/data-structures


*/
func Start(events chan<-) error {
	fmt.Println("starting hose")

	splits := make(chan [8]byte)

	for {
		select {
		case event := <-events:
			// point = event.point() // [8]byte
			// drip_event_channel = find_drip(point)
			// drip_event_channel <- event
		case split := <-splits:
			// close the drip's channel so it can sync and exit
			// mv the old file out of the way and mkdir
			// replace the writer's node in the trie with one that indicates it's a directory
		}
	}
	return nil
}
