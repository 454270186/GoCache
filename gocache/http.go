package gocache

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/454270186/GoCache/gocache/consistenthash"
)

const (
	DefaultBasePath = "/_gocache/"
	DefaultReplicas = 50
)

var _ PeerPicker = (*HTTPPool)(nil)

type HTTPPool struct {
	selfAddr    string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.HashRing // <data_key, peer_addr>
	httpGetters map[string]*httpGetter   // <peer_addr, the_Get()>
}

func NewHTTPPool(addr string) *HTTPPool {
	return &HTTPPool{
		selfAddr:    addr,
		basePath:    DefaultBasePath,
		httpGetters: make(map[string]*httpGetter),
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

// Set() sets the peers of the HTTP pool
func (h *HTTPPool) Set(addrs ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.peers = consistenthash.New(DefaultReplicas, nil)
	h.peers.Add(addrs...)
	for _, addr := range addrs {
		h.httpGetters[addr] = &httpGetter{baseURL: addr + DefaultBasePath}
	}
}

// PickPeer() tries to pick a peer from hash ring according to given key
func (h *HTTPPool) PickPeer(key string) (peer PeerGetter, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if peerAddr := h.peers.Get(key); peerAddr != "" && peerAddr != h.selfAddr {
		h.Log("Pick peer %s", peerAddr)
		return h.httpGetters[peerAddr], true
	}

	return nil, false
}

var _ PeerGetter = (*httpGetter)(nil)

type httpGetter struct {
	baseURL string // the address of remote peer
}

// Get() is the HTTP Client for getting value from remote peer
func (h *httpGetter) Get(group, key string) (string, error) {
	peerURL := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	resp, err := http.Get(peerURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server return %s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading resp body: %v", err)
	}

	return string(bytes), nil
}