package main

import (
	"fmt"

	gcache "github.com/454270186/GoCache/api"
	// "github.com/454270186/GoCache/server"
)

func main() {
	g := gcache.NewGoCache(
		"http://0.0.0.0:8002",
		"http://0.0.0.0:8003",
	)
	g.Put("xiaofei", "123")
	g.Put("feifei", "456")
	val, _ := g.Get("feifei")
	fmt.Println(val)

	// server.CacheServerMain()
}