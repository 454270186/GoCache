package lru

import "gocache/list"

type Cache struct {
	maxBytes int64
	nBytes   int64
	l        *list.List
	cache    map[string]*list.Element

	// Hook on delete a K-V
	onDelete func(key string, val string)
}

func New(maxBytes int64, onDelete func(key string, val string)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		l:        list.New(),
		cache:    make(map[string]*list.Element),
		onDelete: onDelete,
	}
}

func (c *Cache) Get(Key string) (string, bool) {
	if _, ok := c.cache[Key]; !ok {
		return "", false
	}

	ele := c.cache[Key]
	c.l.MoveToTail(ele)
	return ele.Val, true
}

func (c *Cache) Put(key string, val string) {
	if ele, ok := c.cache[key]; ok {
		c.nBytes += int64(len(val) - len(ele.Val))
		ele.Val = val
		c.l.MoveToTail(ele)
	} else {
		newEle := &list.Element{Key: key, Val: val}
		c.l.Add(newEle)
		c.cache[key] = newEle
		c.nBytes += int64(len(newEle.Key) + len(newEle.Val))
	}

	for c.maxBytes > 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	oldest := c.l.GetFirst()
	c.nBytes -= int64(len(oldest.Key) + len(oldest.Val))
	delete(c.cache, oldest.Key)

	c.l.RemoveHead()

	if c.onDelete != nil {
		c.onDelete(oldest.Key, oldest.Val)
	}
}

func (c *Cache) Len() int64 {
	return c.nBytes
}
