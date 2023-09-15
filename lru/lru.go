package lru

import "gocache/list"

type Cache struct {
	maxBytes int64
	nBytes   int64
	l        *list.List
	cache    map[string]*list.Element

	// Hook on delete a K-V
	onDelete func(key string, val list.Value)
}

func New(maxBytes int64, onDelete func(key string, val list.Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		l: list.New(),
		onDelete: onDelete,
	}
}

func (c *Cache) Get(Key string) (list.Value, bool) {
	return nil, true
}

func (c *Cache) Add(key string, val list.Value) {

}

func (c *Cache) RemoveOldest() {

}

func (c *Cache) Len() int {
	return 0
}