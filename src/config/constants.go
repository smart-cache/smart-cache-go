package config

import (
	// "time"
)

const CACHE_SIZE = 20

//const SEED = time.Now().UnixNano()
const SEED = 1

type CacheType int

const (
	LRU            CacheType = 0
	MarkovPrefetch CacheType = 1
	MarkovEviction CacheType = 2
	MarkovBoth     CacheType = 3
)

type DataType string

const DATA_FETCH_TIME = 10