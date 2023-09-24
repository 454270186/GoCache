package main

import (
	"os"
	"strconv"

	"github.com/454270186/GoCache/server"
)

const apiAddr = "127.0.0.1:8080"

func main() {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8001
	}
	isAPI := true

	g := server.InitGroup()

	if isAPI {
		go server.RunAPIServer(apiAddr, g)
	}
	server.RunCacheServer(server.AddrMap[port], server.Addrs, g)
}