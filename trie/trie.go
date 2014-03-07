package trie

import "fmt"

/*
trie of bytes arbitrarily deep

trie struct

trie.PrefixMatch(b []bytes) node


reference:
http://golang.org/doc/play/tree.go



new tree:
 - make a single dir node for the parent directory


tree search(bytes []byte)
   b := bytes[0]
   remainder := bytes[1:]

   if self.is_file?, then recurse:

   return search(remainder)



*/

type Node struct {
	b     byte       // byte to check at this node
	kids  [256]*Node // or nil
	thing string
	// f *io.Writer  // or nil
}

// Match traverses a Node till it finds a FileNode and then returns it
func (dn *Node) Match(data []byte) (*FileNode, error) {
}
