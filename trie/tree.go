package main

import (
	"fmt"
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
