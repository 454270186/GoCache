package gcache

import (
	"fmt"
	"log"

	"github.com/454270186/GoCache/gocache"
)

type Engine struct {
	selfAddr string

	peers     *gocache.HTTPPool
	peerAddrs []string
}

func NewGoCache(peerAddrs ...string) *Engine {
	e := &Engine{}
	baseGroup := gocache.NewGroup("base", 20, func(key string) (string, error) {
		log.Println("随便找找:", key)
		return "sbsbsbsb", nil

	})
	gocache.AddGroup("base", baseGroup)

	if len(peerAddrs) > 0 {
		e.peers = gocache.NewHTTPPool("")
		e.peerAddrs = peerAddrs
		e.peers.Set(peerAddrs...)
		baseGroup.RegisterPeers(e.peers)
	}

	return e
}

func (e *Engine) NewGroup(name string, cacheBytes int64) {
	if g := gocache.GetGroup(name); g != nil {
		return
	}

	newGroup := gocache.NewGroup(name, cacheBytes, nil)
	gocache.AddGroup(name, newGroup)
}

func (e *Engine) Put(key, val string) error {
	baseGroup := gocache.GetGroup("base")
	if baseGroup == nil {
		panic("base group is nil")
	}

	err := baseGroup.Put(key, val)
	if err != nil {
		return err
	}

	log.Printf("[Put] Put <%s -- %s>\n", key, val)
	return nil
}

func (e *Engine) PutWithGroup(groupName, key, val string) error {
	g := gocache.GetGroup(groupName)
	if g == nil {
		return fmt.Errorf("group %s does not exist", groupName)
	}

	err := g.Put(key, val)
	if err != nil {
		return err
	}

	log.Printf("[Put] Put <%s -- %s>\n", key, val)
	return nil
}

func (e *Engine) Get(key string) (val string, ok bool) {
	baseGroup := gocache.GetGroup("base")
	if baseGroup == nil {
		panic("base group is nil")
	}

	val, err := baseGroup.Get(key)
	if err != nil {
		log.Println("[Error]", err)
		return "", false
	}

	return val, true
}

func (e *Engine) GetWithGroup(groupName, key string) (val string, ok bool) {
	g := gocache.GetGroup(groupName)
	if g == nil {
		return "", false
	}

	val, err := g.Get(key)
	if err != nil {
		log.Println("[Error]", err)
		return "", false
	}

	return val, true
}
