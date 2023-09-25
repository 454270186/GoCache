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
	pb "github.com/454270186/GoCache/gocache/gocachepb/gocachepb"
	"google.golang.org/protobuf/proto"
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
	httpPutters map[string]*httpPutter   // <peer_addr, the_Put()>
}

func NewHTTPPool(addr string) *HTTPPool {
	return &HTTPPool{
		selfAddr:    addr,
		basePath:    DefaultBasePath,
		httpGetters: make(map[string]*httpGetter),
		httpPutters: make(map[string]*httpPutter),
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s\n", p.selfAddr, fmt.Sprintf(format, v...))
}

// Get: /<BASEPATH>/<GroupName>/<Key>
// Put: /<BASEPATH>/<GroupName>/<Key>/<Val>
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, DefaultBasePath) {
		http.Error(w, "[Error] bad base path", http.StatusNotFound)
		return
	}

	p.Log("%s %s", r.Method, r.URL.Path)
	httpMethod := r.Method
	switch httpMethod {
	case "GET":
		url := r.URL.Path
		parts := strings.Split(url, "/")
		if len(parts) != 4 {
			http.Error(w, "bad URL", http.StatusBadRequest)
			return
		}

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

		// Decode response into protobuf
		body, err := proto.Marshal(&pb.GetResponse{Value: val})
		if err != nil {
			http.Error(w, fmt.Sprintf("[Error] Decode failed"), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "octet-stream")
		w.Write(body)
	case "POST":
		url := r.URL.Path
		parts := strings.Split(url, "/")
		if len(parts) != 5 {
			http.Error(w, "bad URL", http.StatusBadRequest)
			return
		}

		groupName, key, val := parts[2], parts[3], parts[4]

		group := GetGroup(groupName)
		if group == nil {
			http.Error(w, fmt.Sprintf("[Error] group %s does not exist", groupName), http.StatusNotFound)
			return
		}

		err := group.Put(key, val)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.Log("[Put] Put <%s -- %s>", key, val)
		w.Write([]byte("Successfully put new k-v"))
	}
}

// Set() sets the peers of the HTTP pool
func (h *HTTPPool) Set(addrs ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.peers = consistenthash.New(DefaultReplicas, nil)
	h.peers.Add(addrs...)
	for _, addr := range addrs {
		h.httpGetters[addr] = &httpGetter{baseURL: addr + DefaultBasePath}
		h.httpPutters[addr] = &httpPutter{baseURL: addr + DefaultBasePath}
	}
}

// PickPeer() tries to pick a peer from hash ring according to given key
func (h *HTTPPool) PickPeer(key string) (peerGetter PeerGetter, peerPutter PeerPutter, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if peerAddr := h.peers.Get(key); peerAddr != "" && peerAddr != h.selfAddr {
		h.Log("Pick peer %s", peerAddr)
		return h.httpGetters[peerAddr], h.httpPutters[peerAddr], true
	}

	return nil, nil, false
}

var _ PeerGetter = (*httpGetter)(nil)

type httpGetter struct {
	baseURL string // the address of remote peer
}

// Get() is the HTTP Client for getting value from remote peer
func (h *httpGetter) Get(in *pb.GetRequest, out *pb.GetResponse) error {
	peerURL := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)

	resp, err := http.Get(peerURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Server return %s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading resp body: %v", err)
	}

	if err := proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("error while decoding resp body: %v", err)
	}

	return nil
}

type httpPutter struct {
	baseURL string // the address of remote peer
}

func (h *httpPutter) Put(in *pb.PutRequest, out *pb.PutResponse) error {
	peerURL := fmt.Sprintf(
		"%v%v/%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
		url.QueryEscape(in.GetValue()),
	)

	resp, err := http.Post(peerURL, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Server return %s", resp.Status)
	}

	// bytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return fmt.Errorf("error while reading resp body: %v", err)
	// }

	// if err := proto.Unmarshal(bytes, out); err != nil {
	// 	return fmt.Errorf("error while decoding resp body: %v", err)
	// }

	return nil
}
