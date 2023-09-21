package main

import (
	"flag"

	"github.com/454270186/GoCache/server"
)

const apiAddr = "127.0.0.1:8080"

func main() {
	var port int
	var isAPI bool
	flag.IntVar(&port, "port", 8001, "port of cache server")
	flag.BoolVar(&isAPI, "api", false, "")
	flag.Parse()

	g := server.InitGroup()

	if isAPI {
		go server.RunAPIServer(apiAddr, g)
	}
	server.RunCacheServer(server.AddrMap[port], server.Addrs, g)
}