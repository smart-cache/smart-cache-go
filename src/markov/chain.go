package markov

import (
	"sync"
	"log"
	"math"
	"../heap"
)

type MarkovChain struct {
	nodes			map[string]*MarkovNode  // filename -> Node (with adjacencies)
	lastAccess		map[int]string			// client ID -> lastAccess
	mu				sync.Mutex
}


// creates empty node for the given name
func MakeMarkovChain() *MarkovChain {
	// create empty set of 
	markov := &MarkovChain{
		lastAccess: make(map[int]string), 
		nodes: make(map[string]*MarkovNode),
	}
	// set default MC for first call to markov::Access() (for each client)
	markov.nodes[""] = MakeMarkovNode("")
	return markov
}

func (m *MarkovChain) RecordTransition(filename string, id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	last, ok := m.lastAccess[id]

	if !ok {
		// this client ID's first access
		last = ""
	}

	m.nodes[last].RecordTransition(filename)

	// check if file has own chain
	if _, ok := m.nodes[filename]; !ok {
		m.nodes[filename] = MakeMarkovNode(filename)
	}
	m.lastAccess[id] = filename
}


// predict the next n files after filename is accessed
func (m *MarkovChain) BatchPredict(filename string, n int) []string {
	// this is coarse-gained locking
	m.mu.Lock()
	defer m.mu.Unlock()
	if n < 0 {
		log.Fatalf("AJ hasn't implemented fetching %d files yet :/", n)
	}
	// run Dijkstra's and return the results
	return m.longPaths(filename, n)
}

// Find highest probabilities from source
// assumes source is in m.nodes
// CANNOT predict source as likely to be fetched again
// return order likelihood order
func (m *MarkovChain) longPaths(source string, n int) []string {
	// set up min weights
	distances := make(map[string]float64)
	distances[source] = 0

	// store removed nodes so we can fetch the closest values and check if something has been removed
	removed_nodes := make(map[string]bool)
	closest_files := make([]string, 0)
	nRemoved := 0

	// store current guesses
	var queue heap.MinHeapFloat
	queue.Init()

	// relax all edges from source
	src_node, ok := m.nodes[source]
	if !ok {
		log.Fatalf("THIS SHOULD NEVER HAPPEN [source node not found] %v -> %v", source, m.nodes)
	}

	// initialize with all of the adjacencies of the source node
	for _, neighbor := range src_node.adjacencies {
		// weights are the negated log of the edge ratio -> min path weight becomes max product (max probability)
		weight := -math.Log((float64(neighbor.count) / float64(src_node.count)))
		distances[neighbor.name] = weight
		queue.Insert(neighbor.name, weight)
	}

	// now run Dijkstra's
	for queue.Size > 0 && nRemoved < n {
		name := queue.ExtractMin()
		node := m.nodes[name]
		estimate := distances[name]
		// this file is close in probability, so remove it from valid candidates
		removed_nodes[name] = true
		closest_files = append(closest_files, name)

		// iterate through all neighbors of this file
		for _, transition := range node.adjacencies {
			// check if neighbor file has been seen before
			if _, ok := distances[transition.name]; !ok {
				// not seen before, set probability estimate and insert into heap
				distances[transition.name] = math.Inf(1)
				queue.Insert(transition.name, math.Inf(1))
			}

			if _, ok := removed_nodes[name]; (!ok && transition.name != source) {
				// this neighbor has not been removed already and is not the source node
				// then try to relax weight estimate
				weight := -math.Log((float64(transition.count) / float64(node.count)))
				if (weight + estimate) < distances[transition.name] {
					// then relax this edge
					distances[transition.name] = (weight + estimate)
					// this will insert if not already found
					queue.ChangeKey(transition.name, (weight + estimate))
				}
			}
		}
	}
	return closest_files
}