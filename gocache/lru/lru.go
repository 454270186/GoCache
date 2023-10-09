package lru

import (
	"fmt"
	"time"

	"github.com/454270186/GoCache/gocache/list"
)

type Cache struct {
	maxBytes int64
	nBytes   int64

	old   *list.List
	young *list.List

	cache map[string]*list.Element

	// Hook on delete a K-V
	onDelete func(key string, val string)
}

func New(maxBytes int64, onDelete func(key string, val string)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		old:      list.New(),
		young:    list.New(),
		cache:    make(map[string]*list.Element),
		onDelete: onDelete,
	}
}

func (c *Cache) Get(Key string) (string, bool) {
	if ele, ok := c.cache[Key]; !ok {
		if ele.Heat == list.Old {
			fmt.Printf("key %v is in old\n", ele.Key)
		} else {
			fmt.Printf("key %v is in young\n", ele.Key)
		}
		return "", false
	}

	ele := c.cache[Key]
	if ele.Heat == list.Old {
		fmt.Printf("key %v is in old\n", ele.Key)
	} else {
		fmt.Printf("key %v is in young\n", ele.Key)
	}

	// new
	if ele.Heat == list.Old {
		if time.Since(ele.LastVis) > list.BlockOldInterval {
			ele.Heat = list.Young
			c.young.Add(ele)
			c.old.Remove(ele)
		}
		ele.LastVis = time.Now()
	} else if ele.Heat == list.Young {
		c.young.MoveToTail(ele)
	}

	ele.LastVis = time.Now()

	return ele.Val, true
}

func (c *Cache) Put(key string, val string) {
	if ele, ok := c.cache[key]; ok {
		c.nBytes += int64(len(val) - len(ele.Val))
		ele.Val = val

		// new
		if ele.Heat == list.Old {
			if time.Since(ele.LastVis) > list.BlockOldInterval {
				ele.Heat = list.Young
				c.young.Add(ele)
				c.old.Remove(ele)
			}
		} else if ele.Heat == list.Young {
			// TODO
			c.young.MoveToTail(ele)
		}

		ele.LastVis = time.Now()

	} else {
		newEle := &list.Element{
			Key: key,
			Val: val,
			Heat: list.Old,
			LastVis: time.Now(),
		}

		// new
		c.old.Add(newEle)
		c.cache[key] = newEle
		c.nBytes += int64(len(newEle.Key) + len(newEle.Val))
	}

	for c.maxBytes > 0 && c.nBytes > c.maxBytes {
		// new
		if c.old.Len() > 0 {
			c.old.RemoveHead()
		} else if c.young.Len() > 0 {
			c.young.RemoveHead()
		}
	}
}

// func (c *Cache) RemoveOldest() {
// 	oldest := c.l.GetFirst()
// 	c.nBytes -= int64(len(oldest.Key) + len(oldest.Val))
// 	delete(c.cache, oldest.Key)

// 	c.l.RemoveHead()

// 	if c.onDelete != nil {
// 		c.onDelete(oldest.Key, oldest.Val)
// 	}
// }

func (c *Cache) Len() int64 {
	return c.nBytes
}
