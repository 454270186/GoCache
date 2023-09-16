package gocache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const DefaultBasePath = "/_gocache/"

type HTTPPool struct {
	selfAddr string
	basePath string
}

func NewHTTPPool(addr string) *HTTPPool {
	return &HTTPPool{
		selfAddr: addr,
		basePath: DefaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s\n", p.selfAddr, fmt.Sprintf(format, v...))
}

// <BASEPATH>/<GroupName>/<Key>
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, DefaultBasePath) {
		http.Error(w, "[Error] bad base path", http.StatusNotFound)
		return
	}

	p.Log("%s %s", r.Method, r.URL.Path)
	url := r.URL.Path
	parts := strings.Split(url, "/")

	groupName, key := parts[2], parts[3]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, fmt.Sprintf("[Error] group %s does not exist", groupName), http.StatusNotFound)
		return
	}

	val, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "octet-stream")
	w.Write([]byte(val))
}