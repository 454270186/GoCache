package main

import (
	"fmt"
	// "os"
	// "strconv"

	gcache "github.com/454270186/GoCache/api"
	// "github.com/454270186/GoCache/server"
)

const apiAddr = "127.0.0.1:8080"

// func main() {
// 	portStr := os.Getenv("PORT")
// 	port, err := strconv.Atoi(portStr)
// 	if err != nil {
// 		port = 8001
// 	}
// 	isAPI := true

// 	g := server.InitGroup()

// 	if isAPI {
// 		go server.RunAPIServer(apiAddr, g)
// 	}
// 	server.RunCacheServer(server.AddrMap[port], server.Addrs, g)
// }

func main() {
	g := gcache.NewGoCache(
		"http://0.0.0.0:8002",
		"http://0.0.0.0:8003",
	)
	g.Put("xiaofei", "123")
	g.Put("xiaofei", "456")
	val, _ := g.Get("xiaofei")
	fmt.Println(val)
}