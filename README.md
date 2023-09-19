## GoCache
a simple distributed cache

### Features
- 实现LRU缓存淘汰算法，有效利用内存，提高缓存命中率
- 实现一致性哈希环，提高了节点的均匀性和可拓展性
- 多节点间基于HTTP进行通信，实现缓存节点的水平拓展
- 使用singleflight合并高并发的相同请求，防止缓存击穿