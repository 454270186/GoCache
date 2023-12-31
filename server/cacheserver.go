package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/454270186/GoCache/gocache"
)

// use as DB for testing
var db = map[string]string{
	"xiaofei": "100",
	"dafei":   "500",
	"yuerfei": "250",
	"1": "1",
	"123": "123",
	"sb": "sb",
	"shshsh": "shshsh",
}

var AddrMap = map[int]string{
	8001: "http://0.0.0.0:8001",
	8002: "http://0.0.0.0:8002",
	8003: "http://0.0.0.0:8003",
}

const apiAddr = "127.0.0.1:8080"

var (
	ISAPI bool
	PORT int
)

func init() {
	var err error

	portStr := os.Getenv("PORT")
	PORT, err = strconv.Atoi(portStr)
	if err != nil {
		PORT = 8001
	}
	API := os.Getenv("API")
	ISAPI, err = strconv.ParseBool(API)
	if err != nil {
		ISAPI = false
	}
}

func RunCacheServer(serverAddr string, peerAddrs []string, group *gocache.Group) {
	peers := gocache.NewHTTPPool(serverAddr)
	peers.Set(peerAddrs...)
	
	group.RegisterPeers(peers)

	// Async health check
	if ISAPI {
		peers.StartHealthCheck()
	}
	
	log.Printf("[rpc]Cache Server start listening address: %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr[7:], peers); err != nil {
		log.Fatal(err)
	}
}

func InitGroup() *gocache.Group {
	g := gocache.NewGroup("base", 100, func(key string) (string, error) {
		log.Println("[SlowDB] Search key", key)
		if v, ok := db[key]; ok {
			return v, nil
		}

		return "", fmt.Errorf("%s does not exist", key)
	})

	return g
}

func CacheServerMain() {
	var Addrs = []string{
		"http://0.0.0.0:8001",
		"http://0.0.0.0:8002",
		"http://0.0.0.0:8003",
	}

	g := InitGroup()

	if ISAPI {
		go RunAPIServer(apiAddr, g)
	}
	RunCacheServer(AddrMap[PORT], Addrs, g)
}