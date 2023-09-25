package main

import (
	"fmt"

	gcache "github.com/454270186/GoCache/api"
)

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