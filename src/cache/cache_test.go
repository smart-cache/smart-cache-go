package cache

import (
	"fmt"
	// "reflect"
	"strconv"
	"testing"
	"../datastore"
	// "../utils"
	"../config"
)

func TestBasicLRUFail(t *testing.T) {
	fmt.Printf("TestBasicLRUFail ...\n")
	failed := false

	data := datastore.MakeDataStore()

	// add files to datastore
	for j := 0; j < (config.CACHE_SIZE + 1); j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
        data.Make(filename, config.DataType(filename))
	}

	id := 1
	iter := 4 // number of iterations

	// this copies data, so can't adjust later
	cache := MakeCache(id, config.CACHE_SIZE, config.LRU, data)

	for i := 0; i < iter; i++ {
		for j := 0; j < (config.CACHE_SIZE + 1); j++ {
			filename := "fake_" + strconv.Itoa(j) + ".txt"
			_, err := cache.Fetch(filename, id)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	hits, misses, _ := cache.Report()

	expected_misses := (int64(iter) * (config.CACHE_SIZE + 1))
	if hits != 0 || misses != expected_misses {
		t.Errorf("Expected 0 hits and %d misses, got %d hits and %d misses.", expected_misses, hits, misses)
		failed = true
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

func TestBasicLRUSuccess(t *testing.T) {
	fmt.Printf("TestBasicLRUSuccess ...\n")
	failed := false
	data := datastore.MakeDataStore()

	// add files to datastore
	for j := 0; j < config.CACHE_SIZE; j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
        data.Make(filename, config.DataType(filename))
	}

	id := 1
	iter := 4 // number of iterations

	cache := MakeCache(id, config.CACHE_SIZE, config.LRU, data)

	if config.CACHE_SIZE > 100 {
		fmt.Printf("\tignoring, CACHE_SIZE too big\n")
		return
	}

	for i := 0; i < iter; i++ {
		for j := 0; j < config.CACHE_SIZE; j++ {
			filename := "fake_" + strconv.Itoa(j) + ".txt"
			_, err := cache.Fetch(filename, id)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	hits, misses, _ := cache.Report()
	expected_hits := (config.CACHE_SIZE * int64(iter - 1))

	if hits != expected_hits || misses != config.CACHE_SIZE {
		t.Errorf("Expected %d hits and %d misses, got %d hits and %d misses.", expected_hits, config.CACHE_SIZE, hits, misses)
		failed = true
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

func TestBasicMarkovFail(t *testing.T) {
	fmt.Printf("TestBasicMarkovFail ...\n")
	failed := false

	data := datastore.MakeDataStore()

	// add files to datastore
	for j := 0; j < (config.CACHE_SIZE + 1); j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
        data.Make(filename, config.DataType(filename))
	}

	id := 1
	iter := 4 // number of iterations

	// this copies data, so can't adjust later
	cache := MakeCache(id, config.CACHE_SIZE, config.Markov, data)

	for i := 0; i < iter; i++ {
		for j := 0; j < (config.CACHE_SIZE + 1); j++ {
			filename := "fake_" + strconv.Itoa(j) + ".txt"
			_, err := cache.Fetch(filename, id)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	hits, misses, _ := cache.Report()

	expected_misses := (int64(iter) * (config.CACHE_SIZE + 1))
	if hits == 0 || misses >= expected_misses {
		t.Errorf("Expected more than 0 hits and less than %d misses, got %d hits and %d misses.", expected_misses, hits, misses)
		failed = true
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}
