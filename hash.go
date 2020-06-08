package cache

import (
    "math/rand"
    "../config"
)

/************************************************
Hash Function API

* To Client:
    GetCaches(file string, clientID int) []int
        Return which cache(s) a client should talk to for a particular file

* To CacheMaster:
    MakeHash()
        initiliazation of hash object
    GetCachesInGroup(groupID)
        Get the cache ids that are in a particular group

************************************************/

type Hash struct {
	NumGroups       int
	clientIds       []int
	fileGroups      map[string]int // map of file to column group
	cacheIdToGroup  map[int]int
	groupToCacheIDs map[int][]int
	replicaOrder    map[string]map[int][]int // map of file to client id to an ordering amongst replicas to check
	cacheIDs        []int
}


/************************************************************
API Useful to Client
*************************************************************/
func (h *Hash) GetCaches(file string, clientID int) []int {
	return h.replicaOrder[file][clientID]
}


/************************************************************
API Useful to Cache Master
*************************************************************/
func MakeHash(numCaches int, filenames []string, n int, replication int, clients []int) *Hash {
    h := &Hash{}
    h.initializeClientIDs(clients)
    h.NumGroups = numCaches / replication // number of "columns"
    h.fileGroups = makeFileGroups(filenames, n, h.NumGroups, config.SEED)
    h.cacheIdToGroupInit(numCaches, h.NumGroups)
    h.makeCacheOrderings(filenames)
    return h
}

func (h *Hash) GetCachesInGroup(groupID int) []int {
    return h.groupToCacheIDs[groupID]
}

/***********************************************************
API Useful in Testing
***********************************************************/
func (h *Hash) GetFileGroups(groupID int) map[int][]int {
    return h.groupToCacheIDs
}


/*
Internal Usage in hash creation and initialization
*/
func (h *Hash) makeCacheOrderings(filenames []string) {
    replicaOrder := map[string]map[int][]int{} // filename --> client id --> ordering
    for _, file := range filenames {
        replicaOrder[file] = h.getCacheOrdersForFile(file)
    }
    h.replicaOrder = replicaOrder
}

func (h *Hash) initializeClientIDs(clients []int) {
	ids := make([]int, len(clients))
	for i, client := range clients {
		ids[i] = client
	}
	h.clientIds = ids
}

func (h *Hash) getCacheOrdersForFile(file string) map[int][]int {
    caches := h.groupToCacheIDs[h.fileToGroup(file)]
    mapping := map[int][]int{}
    for i, clientId := range h.clientIds {
        shuffled := shuffle(caches, i)
        mapping[clientId] = make([]int, len(shuffled))
        copy(mapping[clientId], shuffled)
    }
    return mapping
}


func (h *Hash) fileToGroup(filename string) int {
    return h.fileGroups[filename]
}

func (h *Hash) cacheIdToGroupInit(numCaches int, numGroups int) {
    mapping := splitAmongstGroups(numCaches, numGroups)
    idToGroup := make(map[int]int)
    for id := 0; id < numCaches; id++ {
        idToGroup[id] = mapping[id]
    }
    groupToID := make(map[int][]int)
    for id, group := range idToGroup {
        groupToID[group] = append(groupToID[group], id)
    }
    h.cacheIdToGroup = idToGroup
    h.groupToCacheIDs = groupToID
}

func splitAmongstGroups(n int, numGroups int) []int {
    // splits n files across numGroups groups evenly
    mapping := make([]int, n)
    minpergroup := n / numGroups

    for i := 0; i < numGroups; i++ {
        for j := 0; j < minpergroup; j++ {
            mapping[minpergroup*i + j] = i
        }
    }
    //maps the overflow if n not divisible by numGroups
    for i := 0; i < n - (numGroups * minpergroup); i++ {
        mapping[minpergroup*numGroups + i] = i
    }
    return mapping

}

func shuffle(slice []int, seed int) []int {
    rand.Seed(int64(seed))
    rand.Shuffle(len(slice), func(i, j int) {
        slice[i], slice[j] = slice[j], slice[i]
    })
    return slice
}

func makeFileGroups(filenames []string, n int, numGroups int, seed int) map[string]int {
    mapping := splitAmongstGroups(n, numGroups)
    mapping = shuffle(mapping, seed)
    fileGroups := make(map[string]int)
    for i := 0; i < n; i++ {
        fileGroups[filenames[i]] = mapping[i]
    }
    return fileGroups
}
