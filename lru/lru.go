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
	return nil
}