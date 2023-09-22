## GoCache
a simple distributed cache

### Features
- 实现LRU缓存淘汰算法，有效利用内存，提高缓存命中率
- 实现一致性哈希环，提高了节点的均匀性和可拓展性
- 多节点间基于HTTP进行通信，实现缓存节点的水平拓展
- 使用singleflight合并高并发的相同请求，防止缓存击穿

### Usage

#### Single node
```go
package main

import (
	gcache "github.com/454270186/GoCache/api"
)

func main() {
	g := gcache.NewGoCache()
	g.Put("xiaofei", "100")
	val, _ := g.Get("xiaofei") // val ==> 100
}
```

You can also store the data in specific group

```go
package main

import (
	gcache "github.com/454270186/GoCache/api"
)

func main() {
	g := gcache.NewGoCache()

	g.NewGroup("people", 20)
	g.PutWithGroup("people", "xiaofei", "1")
	g.NewGroup("student", 20)
	g.PutWithGroup("student", "xiaofei", "100")

	val1, _ := g.GetWithGroup("people", "xiaofei")  // val1 ==> 1
	val2, _ := g.GetWithGroup("student", "xiaofei") // val2 ==> 100
}
```

PS: If dont specify a group, the data will store in "base" group