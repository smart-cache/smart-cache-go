package cache

import (
	"sync"
	"log"
	"errors"

	"../heap"
	"../markov"
	"../datastore"
	"../config"
)

/********************************
Cache supports the following external API to users
c.Init(cacheSize int, cacheType CacheType, data *datastore.DataStore)
	Initializes a cache with eviction policy and prefetch defined by cache type
	Copies underlying datastore
c.Report() (hits, misses)
    Get a report of the hits and misses  TODO: Do we want a version number or
    timestamp mechanism of any form here?
c.Fetch(name string) (config.DataType, error)
*********************************/
type Cache struct {
	mu          sync.Mutex          			// Lock to protect shared access to cache
	cache	    map[string]config.DataType		// cached data storage
	heap		*heap.MinHeapInt64				// for LRU version
	timestamp	int64 							// for controlling LRU heap
	maxSize		int64							// maximum allowable cache size
	chain		*markov.MarkovChain				// for Markov version
	cType		config.CacheType
	data		*datastore.DataStore			// for fetching data

	// external data
	id          int								// uid for each cache (provided by ctor)
	misses		int64
	hits		int64
}

// creates a copy by copying the underlying datastore
func MakeCache(id int, cacheSize int64, cacheType config.CacheType, data *datastore.DataStore) (* Cache) {
	cache := &Cache{
		// set user provided vars
		cType: cacheType,
		id: id,
		maxSize: cacheSize,
		data: data.Copy(),

		// set type defined vars
		misses: 0,
		hits: 0,
		cache: make(map[string]config.DataType),
		timestamp: 0,

		// set special datatypes
		heap: heap.MakeMinHeapInt64(),
		chain: markov.MakeMarkovChain(),
	}
	return cache
}

func (cache *Cache) Fetch(filename string) (config.DataType, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	file, ok := cache.cache[filename]
	cache.timestamp++

	// inform the markov chain of this transaction
	cache.chain.RecordTransition(filename, cache.id)
	// and inform the heap
	cache.heap.ChangeKey(filename, cache.timestamp)

	var err error

	if ok {
		cache.hits++
		err = nil
	} else {
		file, err = cache.AddFileToCache(filename)
		cache.misses++
	}

	// TODO: may want to change the ordering of the prefetching
	if cache.timestamp % config.PREFETCH_SIZE == 0 {
		go cache.BatchPrefetch(filename)
	}
	return file, err
}

func (cache *Cache) Report() (int64, int64, int64) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	return cache.hits, cache.misses, cache.data.CountCalls()
}

func (cache *Cache) BatchPrefetch (filename string) {
	if cache.cType != config.LRU {
		files := cache.chain.BatchPredict(filename, config.PREFETCH_SIZE)
		cache.mu.Lock()
		cache.AddBatchToCache(files)
		cache.mu.Unlock()
	}
}

// assumes lock on cache.mu is held
func (cache *Cache) AddFileToCache(filename string) (config.DataType, error) {
	file, ok := cache.cache[filename]

	if !ok {
		file, ok = cache.data.Get(filename)

		if !ok {
			log.Fatalf("Failed to fetch file %v from underlying datastore", filename)
		}

		// fill the cache with this new datatype
		cache.AddFile(filename, file)
	}

	var err error

	if ok {
		err = nil
	} else {
		err = errors.New("Failed to fetch file from cache")
	}
	return file, err
}

// assumes lock on cache.mu is held
func (cache *Cache) AddFile(filename string, file config.DataType) {
	cache.cache[filename] = file
	cache.heap.Insert(filename, cache.timestamp)

	if cache.heap.Size > cache.maxSize {
		// need to evict, so remove least recently used item
		evict := cache.heap.ExtractMin()
		delete(cache.cache, evict)
		if cache.heap.Size > cache.maxSize {
			log.Fatalf("Cache eviction did not properly fix size: %v > %v", cache.heap.Size, cache.maxSize)
		}
	}
}

// assumes lock on cache.mu is held
func (cache *Cache) AddBatchToCache(filenames []string) (error) {

	files, ok := cache.data.GetBatch(filenames)

	if !ok {
		log.Fatalf("Failed to fetch batch <%v> from underlying datastore", filenames)
	}

	for i, filename := range filenames {
		cache.AddFile(filename, files[i])
	}

	return nil
}