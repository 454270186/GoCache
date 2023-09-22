package main

import (
	// "flag"

	"fmt"

	gcache "github.com/454270186/GoCache/api"
	// "github.com/454270186/GoCache/server"
)

const apiAddr = "127.0.0.1:8080"

// func main() {
// 	var port int
// 	var isAPI bool
// 	flag.IntVar(&port, "port", 8001, "port of cache server")
// 	flag.BoolVar(&isAPI, "api", false, "")
// 	flag.Parse()

// 	g := server.InitGroup()

// 	if isAPI {
// 		go server.RunAPIServer(apiAddr, g)
// 	}
// 	server.RunCacheServer(server.AddrMap[port], server.Addrs, g)
// }

func main() {
	g := gcache.NewGoCache()
	g.Put("xiaofei", "100")
	val, _ := g.Get("xiaofei")
	fmt.Println(val)
}