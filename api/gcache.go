package gcache

import (
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
