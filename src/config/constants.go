package config

import (
	"time"
)

const CACHE_SIZE = 20
const PREFETCH_SIZE = 10

//const SEED = time.Now().UnixNano()
const SEED = 1

type CacheType int

const (
	LRU            	CacheType = 0
	Markov			CacheType = 1
)

type DataType string

const DATA_FETCH_TIME = time.Millisecond * 10
const DATA_COST_TIME = time.Millisecond * 1