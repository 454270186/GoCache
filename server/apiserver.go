package server

import (
	"log"
	"net/http"

	"github.com/454270186/GoCache/gocache"
)

// http://host:port/api?key=xxx
func RunAPIServer(addr string, group *gocache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		v, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "octet-stream")
		w.Write([]byte(v))
	}))

	log.Printf("API Server start listening at %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}