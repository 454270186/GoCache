package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"time"
)

const (
	HealthCheckInterval = 5 * time.Second
)

type Hash func(data []byte) uint32

type HashRing struct {
	hashfunc Hash
	replicas int
	keys     []int
	hashMap  map[int]string

	RealKeysMap map[int]bool
}

func New(nReplicas int, hashfn Hash) *HashRing {
	hr := &HashRing{
		replicas: nReplicas,
		hashMap:  make(map[int]string),
		keys:     make([]int, 0),
		RealKeysMap: make(map[int]bool),
		hashfunc: hashfn,
	}

	if hashfn == nil {
		hr.hashfunc = crc32.ChecksumIEEE
	}

	return hr
}

// Add peer-nodes into hash ring
func (hr *HashRing) Add(keys ...string) {
	for _, key := range keys {
		var hash int
		for i := 0; i < hr.replicas; i++ {
			hash = int(hr.hashfunc([]byte(strconv.Itoa(i) + key)))
			hr.hashMap[hash] = key
			hr.keys = append(hr.keys, hash)
		}

		hr.RealKeysMap[hash] = true
	}

	sort.Ints(hr.keys)
}

// Remove a peer npde and its replicas
func (hr *HashRing) Remove(key string) {
	keysToRemove := make(map[int]bool)

	for k, v := range hr.hashMap {
		if v == key {
			if hr.RealKeysMap[k] {
				delete(hr.RealKeysMap, k)
			}

			keysToRemove[k] = true
			delete(hr.hashMap, k)
		}
	}

	newKeys := make([]int, 0, len(hr.keys))
	for _, k := range hr.keys {
		if !keysToRemove[k] {
			newKeys = append(newKeys, k)
		}
	}

	hr.keys = newKeys
}

// Get hash node by given data key
func (hr *HashRing) Get(key string) string {
	if len(hr.keys) == 0 {
		return ""
	}

	hash := int(hr.hashfunc([]byte(key)))
	idx := sort.Search(len(hr.keys), func(i int) bool {
		return hr.keys[i] >= hash
	})

	return hr.hashMap[hr.keys[idx%len(hr.keys)]]
}