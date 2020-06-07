package markov

import (
	"testing"
	"fmt"
)

func MakeAccesses(m *MarkovChain, files []string, id int) {
	for _, file := range files {
		m.RecordTransition(file, id)
	}
}

func CheckPredictions(received []string, expected []string, t *testing.T) bool {
	if len(received) != len(expected) {
		t.Errorf("CheckPredictions FAILED. Received: size %v instead of size %v", len(received), len(expected))
	}
	for index, file := range received {
		if (file != expected[index]) {
			t.Errorf("CheckPredictions FAILED. Received: file %v instead of file %v", file, expected[index])
			return false
		}
	}
	return true
}

func TestChainSimple(t *testing.T) {
	fmt.Printf("TestChainSimple ...\n")
	failed := false

	chain := MakeMarkovChain()

	files := []string{
		"a.png",
		"b.png",
		"c.png",
		"a.png",
		"b.png",
		"c.png",
	}

	expected_predict := []string{
		"b.png",
		"c.png",
	}

	// execute transitions
	MakeAccesses(chain, files, 1)

	// first try predicting the next two values
	short_predict := chain.BatchPredict("a.png", 2)

	if !CheckPredictions(short_predict, expected_predict, t) {
		t.Errorf("short_predict failed")
	}

	// still expect it to only find the two files
	long_predict := chain.BatchPredict("a.png", 20)
	
	if !CheckPredictions(long_predict, expected_predict, t) {
		t.Errorf("long_predict failed")
	}


	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}