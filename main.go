package main

import (
	"fmt"
	"gocache/lru"
)

func main() {
	c := lru.New(10, nil)

	c.Put("1", "1")
	c.Put("2", "2")
	c.Put("3", "3")
	c.Put("4", "4")
	c.Put("5", "5")
	c.Get("1")
	c.Put("6", "6")

	fmt.Println(c.Len())
	fmt.Println(c.Get("1"))
}