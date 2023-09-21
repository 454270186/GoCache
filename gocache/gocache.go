package gocache

import (
	"errors"
	"log"
	"sync"

	"github.com/454270186/GoCache/gocache/singleflight"
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
	loader    *singleflight.Group
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
		loader: &singleflight.Group{},
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

func (g *Group) Name() string {
	return g.name
}

func (g *Group) Put(key, val string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	return nil
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
	v, err := g.loader.Do(key, func() (interface{}, error) {
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
	})
	if err != nil {
		return "", err
	}

	return v.(string), nil
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

//Add K-V

// populateCache() add a K-V in local cache
func (g *Group) populateCache(key, val string) {
	g.mainCache.add(key, val)
}

func (g *Group) populatePeerCache(key, val string) {

}