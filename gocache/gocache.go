package gocache

import (
	"errors"
	"log"
	"sync"
)

// Callback Func for load data from remote source
type GetterFunc func(key string) (string, error)

/*
Group
*/
type Group struct {
	name      string
	getter    GetterFunc
	mainCache cache
	peers     PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter GetterFunc) *Group {
	if getter == nil {
		panic("Getter func is nil")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("key cannot be empty")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GoCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("Call RegisterPeers() more than once")
	}

	g.peers = peers
}

// load() will try to get data from other data source;
// First try to get from peer cache;
// Finally get data locally;
func (g *Group) load(key string) (string, error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			v, err := peer.Get(g.name, key)
			if err == nil {
				return v, nil
			}
			log.Println("error while get data from peer:", err.Error())
		}
	}

	return g.getLocally(key)
}

// getLocally() calls the Getter callback func to get data
func (g *Group) getLocally(key string) (string, error) {
	v, err := g.getter(key)
	if err != nil {
		return "", err
	}

	g.populateCache(key, v)
	return v, nil
}

func (g *Group) populateCache(key, val string) {
	g.mainCache.add(key, val)
}