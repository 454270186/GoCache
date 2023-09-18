package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type HashRing struct {
	hashfunc Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(nReplicas int, hashfn Hash) *HashRing {
	hr := &HashRing{
		replicas: nReplicas,
		hashMap:  make(map[int]string),
		keys:     make([]int, 0),
		hashfunc: hashfn,
	}

	if hashfn == nil {
		hr.hashfunc = crc32.ChecksumIEEE
	}

	return hr
}

func (hr *HashRing) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < hr.replicas; i++ {
			hash := int(hr.hashfunc([]byte(strconv.Itoa(i) + key)))
			hr.hashMap[hash] = key
			hr.keys = append(hr.keys, hash)
		}
	}

	sort.Ints(hr.keys)
}

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