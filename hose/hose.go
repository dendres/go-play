package hose

import (
	"fmt"
	"github.com/dendres/go-play/event"
)

// the static entries count must match the encoding of each character in path
// currently base32, so it can be 32
// would be cooler to have separate types for leaf and branch nodes
// but I can't figure out how to pass around and use the interface correctly
type Node struct {
	children    [32]*Node
	parent      *Node
	Byte        byte
	HasChildren bool
	IsRoot      bool
}

func (node *Node) Child(b byte) *Node {
	c := node.children[b]
	if c != nil {
		return c
	}
	return node
}

func (node *Node) ReplaceChild(id byte, new_child *Node) {
	node.children[id] = new_child
}

// Search traverses a byte tree and returns the Node representing the closest match
// also returns the unmatched remainder of the id
// does not create or delete nodes
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

	if next_node.HasChildren {
		fmt.Println("recursing to the next matching node")
		return next_node.Search(remainder)
	}
	fmt.Println("returning the next node")
	return next_node, remainder
}

// func (node *Node) Delete() {
//     get parent
//     replace entry in parent with new empty node???
// }

func main() {
	root := Node{IsRoot: true, HasChildren: true}
	id := []byte{0x01, 0x02, 0x03}

	n1, r1 := root.Search(id)
	fmt.Println("result id =", id, "found node =", n1, "remainder =", r1)

	byte1 := byte(0x01)
	node1 := Node{Byte: byte(0x01)}

	root.ReplaceChild(byte1, &node1)

	n1, r1 = root.Search(id)
	fmt.Println("result id =", id, "found node =", n1, "remainder =", r1)

}

type Hose struct {
	root   *Node
	events chan *EventBytes // events to route
	splits chan *Node       // the node to split
}

// search starting at the root node
// find the most specific drip available
// send the event on that channel
func (hose *Hose) Route(eb *EventBytes) {
	point := eb.PointBytes()
	node := hose.root.Search(point)
	// determine if there is a drip available here!!!?????
	// make the Drip type!!!!
	node.Channel <- event
}

// close the drip's channel so it can sync and exit
// mv the old file out of the way and mkdir
// replace the writer's node with one that indicates it's a directory
func (hose *Hose) Split(point []byte) {

}

/*
keeps a prefix tree of channels to drip goroutines
sends each incoming event to the correct drip

watches a split channel that any writer can send to.
pauses event processing and updates the prefix tree when a split is received

*/
func (hose *Hose) Run() error {
	fmt.Println("starting hose")
	for {
		// XXX implement priority here with case high, then case high med, then case high med low
		select {
		case event := <-hose.events:
			hose.Route(event)
		case point_to_split := <-splits:
			hose.Split(point_to_split)
		}
	}
	return nil
}
