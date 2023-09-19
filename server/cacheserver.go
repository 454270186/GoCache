package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/454270186/GoCache/gocache"
)

// use as DB for testing
var db = map[string]string{
	// "xiaofei": "100",
	"dafei":   "500",
	"yuerfei": "250",
}

var AddrMap = map[int]string{
	8001: "http://127.0.0.1:8001",
	8002: "http://127.0.0.1:8002",
	8003: "http://127.0.0.1:8003",
}

var Addrs = []string{
	"http://127.0.0.1:8001",
	"http://127.0.0.1:8002",
	"http://127.0.0.1:8003",
}

func RunCacheServer(serverAddr string, peerAddrs []string, group *gocache.Group) {
	peers := gocache.NewHTTPPool(serverAddr)
	peers.Set(peerAddrs...)
	
	group.RegisterPeers(peers)

	log.Printf("Cache Server start listening address: %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr[7:], peers); err != nil {
		log.Fatal(err)
	}
}

func InitGroup() *gocache.Group {
	g := gocache.NewGroup("student", 100, func(key string) (string, error) {
		log.Println("[SlowDB] Search key", key)
		if v, ok := db[key]; ok {
			return v, nil
		}

		return "", fmt.Errorf("%s does not exist", key)
	})

	return g
}