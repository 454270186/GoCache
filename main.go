package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/454270186/GoCache/gocache"
)

const addr = "127.0.0.1:8080"

// use as DB for testing
var db = map[string]string{
	"xiaofei": "100",
	"dafei":   "500",
	"yuerfei": "250",
}

func main() {
	gocache.NewGroup("student", 100, func(key string) (string, error) {
		log.Println("[SlowDB] Search key", key)
		if v, ok := db[key]; ok {
			return v, nil
		}

		return "", fmt.Errorf("%s does not exist", key)
	})

	peers := gocache.NewHTTPPool(addr)

	log.Printf("Start listening address: %s\n", addr)
	if err := http.ListenAndServe(addr, peers); err != nil {
		log.Fatal(err)
	}
}
