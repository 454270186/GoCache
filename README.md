## GoCache
a simple distributed cache

### Features
- 实现LRU缓存淘汰算法，并对数据进行冷热分区，提高缓存命中率
- 实现一致性哈希环，提高了节点的均匀性和可拓展性
- 多节点间基于HTTP, 使用protobuf进行通信
- 定时心跳检测，及时移除不可用的节点
- 使用singleflight合并高并发的相同请求，防止缓存击穿
- 支持docker水平部署节点

### Getting GoCache
```bash
docker pull erfeiyu/go-cache:latest

docker run -e PORT=8002 -p 8002:8002 erfeiyu/go-cache:latest  # run a cache node in http://0.0.0.0:8002
```

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


#### Multi nodes
```go
package main

import (
	"fmt"

	gcache "github.com/454270186/GoCache/api"
)

func main() {
	g := gcache.NewGoCache(
		"http://0.0.0.0:8002",
		"http://0.0.0.0:8003",
	)
	g.Put("xiaofei", "123")
	g.Put("dafei", "456")
	val, _ := g.Get("xiaofei") // val ==> 123
}
```

**Output**

PeerPicker works
```
2023/09/24 11:43:49 [Server ] Pick peer http://0.0.0.0:8002
2023/09/24 11:43:49 [Put] Put <xiaofei -- 123>
2023/09/24 11:43:49 [Server ] Pick peer http://0.0.0.0:8003
2023/09/24 11:43:49 [Put] Put <dafei -- 456>
2023/09/24 11:43:49 [Server ] Pick peer http://0.0.0.0:8002
123
```
