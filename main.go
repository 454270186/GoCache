package main

import (
	"log"
	"net/http"

	"github.com/454270186/GoCache/gocache"
)

const addr = "127.0.0.1:8080"

func main() {
	peers := gocache.NewHTTPPool(addr)
	
	log.Printf("Start listening address: %s\n", addr)
	if err := http.ListenAndServe(addr, peers); err != nil {
		log.Fatal(err)
	}
}