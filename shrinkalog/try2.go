package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

/*
for compression only... ignore search for now
find the longest repeating []byte in the file

XXX later, events will be stored with delimiter of some kind. newline or offset delimited
    at that time, each event's size will be available. and a buffer can have the whole event
    so won't have to worry about things getting cut off by the buffering
  so just read in the whole file for now as one iteration

http://en.wikipedia.org/wiki/Longest_repeated_substring_problem
make suffix tree, then find the highest node with at least 2 descendants

http://www.allisons.org/ll/AlgDS/Tree/Suffix/
http://en.wikipedia.org/wiki/Suffix_tree
http://en.wikipedia.org/wiki/Suffix_array

http://stackoverflow.com/questions/9452701/ukkonens-suffix-tree-algorithm-in-plain-english

https://www.cs.helsinki.fi/u/ukkonen/SuffixT1withFigs.pdf


https://www.youtube.com/watch?v=F3nbY3hIDLQ
suffix tree:
 - if order doesn't matter, node = hash map
 - if order does matter, then node = "tray" ????


*/

type Un struct {
	From  uint
	To    uint
	Child *Un
}


// https://github.com/kvh/Python-Suffix-Tree/blob/master/suffix_tree.py
// https://code.google.com/p/suffixtree/source/checkout
// https://github.com/JDonner/SuffixTree
// http://felix-halim.net/pg/suffix-tree/



// Note that r.addTransition(...) adds an edge from state r, labelling the edge with a substring. New txt[i]-transitions must be "open" transitions of the form (L,∞).
// ??????

// (s, (k, i-1)) is the canonical reference pair for the active point
func upDate(s, k, i) {
	var oldr = root
	var (endPoint, r) = test_and_split(s, k, i-1, Txt.charAt(i));

	while (!endPoint)
    { r.addTransition(i, infinity, new State());
		if (oldr != root) oldr.sLink = r; // build suffix-link active-path

		oldr = r;
		var (s,k) = canonize(s.sLink, k, i-1)
		(endPoint, r) = test_and_split(s, k, i-1, Txt.charAt(i))
    }

	if(oldr != root) oldr.sLink = s;

	return new pair(s, k);
}//upDate


function test_and_split(s, k, p, t)
{ if(k<=p)
    { // find the t_k transition g'(s,(k',p'))=s' from s
      // k1 is k'  p1 is p' in Ukkonen '95
		var ((k1,p1), s1)  = s[Txt.charAt(k)];

		if (t == Txt.charAt(k1 + p - k + 1))
		return new pair(true, s);
      else
		{ var r = new State()
			s.addTransition(k1, k1+p-k,   r);     // s---->r---->s1
			r.addTransition(    k1+p-k+1, p1, s1);
			return new pair(false, r)
		}
    }
   else // k > p;  ? is there a t-transition from s ?
	return new pair(s[t] != null, s);
}//test_and_split


function canonize(s, k, p)    // s--->...
{ if(p < k) return new pair (s, k);

	// find the t_k transition g'(s,(k',p'))=s' from s
	// k1 is k',  p1 is p' in Ukk' '95
	var ((k1,p1), s1) = s[Txt.charAt(k)];     // s--(k1,p1)-->s1

	while(p1-k1 <= p-k)                       // s--(k1,p1)-->s1--->...
    { k += p1 - k1 + 1;  // remove |(k1,p1)| chars from front of (k,p)
		s = s1;
		if(k <= p)
		{ ((k1,p1), s1) = s[Txt.charAt(k)];   // s--(k1,p1)-->s1
		}
    }
	return new pair(s, k);
}//canonize


function ukkonen95()// construct suffix tree for Txt[0..N-1]
{ var s, k, i;
   var bt;

	root = new State();
	bt = new State();                            // bt (bottom or _|_)

   // Want to create transitions for all possible chars
   // from bt to root
	for (i=0; i < Txt.length; i++)
	bt.addTransition(i,i, root);

	root.sLink = bt;
	s=root; k=0;    // NB. k=0, unlike Ukkonen our strings are 0 based

	for(i=0; i < Txt.length; i++)
    { var (s,k) = upDate(s, k, i);   // follow path from active-point
		(s,k) = canonize(s, k, i);
    }
}//ukkonen95







// active point, which is a triple (active_node,active_edge,active_length)

func main() {

	buf, err := ioutil.ReadFile("test1.log")
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}

	root := Un{}
	current := uint(0)
	// remainder := uint(0) ????

	for offset, runeValue := range buf {
		current = offset
		fmt.Println(offset, string(runeValue))
		log.Fatal("finished")
	}

	// implement some solution to the longest repeated substring problem

}
