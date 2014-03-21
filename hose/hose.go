/*
package hose routes incoming events to the correct file writing "drip".
It keeps a prefix tree of channels to drips.
It performs a "split" when a drip file becomes too large.

XXX splits should mean replacing a file with a directory

1 byte / 2 hex characters per file/directory:
bucket/00/11/22/33/44/55/66/77/00/11/22/33

XXX replace "point" with "point + crc"


*/
package hose

import (
	"fmt"
	"github.com/dendres/go-play/drip"
	"github.com/dendres/go-play/event"
)

// A Node is a branch/leaf in the hose prefix tree
// the static entries count must match the encoding of each character in path
// currently base32, so it can be 32
// would be cooler to have separate types for leaf and branch nodes
// but I can't figure out how to pass around and use the interface correctly
type Node struct {

	// Send EventBytes to the Drip on this channel
	// the channel will block if the Drip is splitting????
	Events chan<- *EventBytes

	// the full path to the file or directory
	Path []byte

	// the byte is the the file or directory name as int8 == byte
	Byte byte

	// true if this is the root node in the tree
	IsRoot bool

	// same as is_directory
	HasChildren bool

	// pointer to the parent Node
	parent   *Node
	children [32]*Node
}

func (node *Node) Child(b byte) *Node {
	c := node.children[b]
	if c != nil {
		return c
	}
	return node
}

// Search returns the Node where the given id should be written
// also returns the unmatched remainder of the id
// if the required node is not an open file, the file is created if needed and then opened for writing
// XXX pass in the channels
func (node *Node) Search(id []byte) (*Node, []byte) {
	fmt.Println("searching id =", id, "node.Byte =", node.Byte)

	nextbyte, remainder := id[0], id[1:]

	if len(remainder) < 1 {
		fmt.Println("ran out of id")
		return node, id
	}

	next_node := node.Child(nextbyte)
	if next_node.Byte != nextbyte {
		fmt.Println("found child with wrong byte: next_node.Byte =", next_node.Byte, ", nextbyte =", nextbyte)
		return node, id
	}

	// XX replace HasChildren with a check to see if the drip is in a state ready to receive
	// determine if there is a drip available here!!!?????
	if next_node.HasChildren {
		fmt.Println("recursing to the next matching node")
		return next_node.Search(remainder)
	}

	fmt.Println("returning the next node")
	return next_node, remainder
}

// ReplaceChild writes the given node to the given id.
func (node *Node) ReplaceChild(id byte, new_child *Node) {
	node.children[id] = new_child
}

// Delete replaces the entry in the parent node with nil.
// XXX must test for memory leek!!!
func (node *Node) Delete() {
	node.parent.children[node.Byte] = nil
}

/*
A Hose keeps a prefix tree of channels to drip goroutines.
It sends each incoming event to the correct drip.
Drips write to the splits channel when they are big enough to split.
Hose performs the split operation creating new drips.
*/
type Hose struct {

	// the root node of the drip tree
	root *Node

	// the incoming channel of events to route
	Events <-chan *EventBytes

	// Receive State changes from the Drips on this channel
	State <-chan string
}

// search starting at the root node
// find the most specific drip available
// send the event on that channel
// Route must NOT be called if any split is in progress!!
func (hose *Hose) Route(eb *EventBytes) {
	point := eb.PointBytes()

	node := hose.root.Search(point)

	// need to create child?????

	// make the Drip type!!!!
	node.Channel <- event
}

// close the drip's channel so it can sync and exit
// mv the old file out of the way and mkdir
// replace the writer's node with one that indicates it's a directory
func (hose *Hose) Split(point []byte) {
	// find the node... like route
	// split the node?????
}

func (hose *Hose) Run() error {
	fmt.Println("starting hose")
	for {
		// XXX implement priority here with case high, then case high med, then case high med low
		select {
		case event := <-hose.Events:
			hose.Route(event)
		case command := <-hose.State:
			select {
			case command == 'split':
				// XXXXXXXXXXXXXXXXXXXX do stuff based on commands from the client... or maybe use separate channels for each command???
			case command == 'split':
			case command == 'split':
			}

		}
	}
	return nil
}

// func main() {
// 	root := Node{IsRoot: true, HasChildren: true}
// 	id := []byte{0x01, 0x02, 0x03}

// 	n1, r1 := root.Search(id)
// 	fmt.Println("result id =", id, "found node =", n1, "remainder =", r1)

// 	byte1 := byte(0x01)
// 	node1 := Node{Byte: byte(0x01)}

// 	root.ReplaceChild(byte1, &node1)

// 	n1, r1 = root.Search(id)
// 	fmt.Println("result id =", id, "found node =", n1, "remainder =", r1)
// }
