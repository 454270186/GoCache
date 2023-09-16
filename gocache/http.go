package gocache

import (
	"fmt"
	"log"
	"net/http"
)

const defaultBasePath = "/_gocache"

type HTTPPool struct {
	selfAddr string
	basePath string
}

func NewHTTPPool(addr string) *HTTPPool {
	return &HTTPPool{
		selfAddr: addr,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s\n", p.selfAddr, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Request] %s %s\n", r.Method, r.URL.Path)
}