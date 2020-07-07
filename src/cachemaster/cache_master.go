package cache_master

import (
	"sync"
	// "time"
	"../datastore"
	"../markov"
	"../config"
	"../cache"
)

/************************************************
Cache Master API
Initialization:
    m = StartTask(
            clientIds       []int
            cacheType   CacheType - specification for prefetch and eviction policies
            numCaches         int - number of cache machines to use
            replication       int - replication factor
            datastore   Datastore
        )
    Initialize a cache master with client list, and replication factor (r)
syncCaches
*************************************************/

type CacheMaster struct {
	mu			sync.Mutex						// lock on master structure
	clientIDs	[]int							// list of all client IDs (TODO: rm if unnecessary)
	caches		map[int]*cache.Cache			// map of cache ID -> cache
	cacheType	config.CacheType				// cache type of all caches (TODO: rm if unnecessary)
	rFactor		int 							// replication factor
	nCaches		int 							// number of caches
	nFiles		int 							// number of pieces of data	(TODO: rm if unnecessary)
	datastore	*datastore.DataStore			// underlying datastore that all caches have access to (TODO: rm if redundant)
	hash		*Hash							// underlying hash method for splitting data access across caches
	sync_time	int 							// how often caches are synced
	chain		*markov.MarkovChain				// most recent aggregate data from syncing
}

type CacheParams struct {
	NCaches 		int 						// number of caches
	RFactor 		int							// replication factor
	CacheType 		config.CacheType			// which type of cache to use (LRU | Markov)
	CacheSize 		int							// size of each cache (assumes homogeneity)
	Datastore 		*datastore.DataStore		// underlying datastore that all caches have access to (TODO: should it be designed this way?)
	Sync_ms 		int							// how many milliseconds to wait in between cache syncs 
}

func MakeCacheMaster(clientIDs []int, params CacheParams) (* CacheMaster) {
	// k: number of caches
	// r: replication factor for data desired
	// this is trivial (can store everything) if cacheSize >= nr/k (where n is
	// size of datastore)
	cm := &CacheMaster{
		clientIDs: clientIDs,
		nCaches: params.NCaches,
		rFactor: params.RFactor,
		datastore: params.Datastore,
		nFiles: params.Datastore.Size(),
		chain: markov.MakeMarkovChain(),
		sync_time: params.Sync_ms,
		caches: make(map[int]*cache.Cache),
	}

	for i := 0; i < cm.nCaches; i++ {
		// datastore is copied in cache making
		c := cache.MakeCache(i, int64(params.CacheSize), params.CacheType, params.Datastore)
		cm.caches[i] = c
	}


	cm.hash = MakeHash(cm.nCaches, cm.datastore.GetFileNames(), cm.nFiles, cm.rFactor, cm.clientIDs)

    if (params.CacheType != config.LRU) {
		// TODO: add periodic syncing of caches
        // go cm.syncCaches(params.Sync_ms)
    }

	return cm
}