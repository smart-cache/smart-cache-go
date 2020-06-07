package markov

import (
	"sync"
)

// Individual edge in the markov graph
// count represents frequency
type MarkovEdge struct {
	count			int
	name			string
}

// sparse representation of adjacencies. double space for efficient lookups + iteration
type MarkovNode struct {
	name			string
	count			int
	adjacencies		[]MarkovEdge			// fast iterator 
	neighbors		map[string]int			// filename -> index in adjacencies. fast lookup of edge weights
	mu				sync.Mutex				// for concurrent requests
}


// creates empty node for the given name
func MakeMarkovNode(name string) *MarkovNode {
	node := &MarkovNode{
		name: name, 
		count: 0, 
		adjacencies: make([]MarkovEdge, 0), 
		neighbors: make(map[string]int),
	}
	return node
}

func (mn *MarkovNode) RecordTransition(filename string) {
	mn.mu.Lock()
	defer mn.mu.Unlock()

	// increase total number of transitions
	mn.count++

	neighbor, ok := mn.neighbors[filename]

	if ok {
		// already have edge to this node
		mn.adjacencies[neighbor].count++
	} else {
		// don't have edge, must make one
		var e MarkovEdge
		e.count = 1 		// first time seeing this transition
		e.name = filename

		// set index in map and append to end of list
		mn.neighbors[filename] = len(mn.adjacencies)
		mn.adjacencies = append(mn.adjacencies, e)
	}
}