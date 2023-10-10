package gocache

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/454270186/GoCache/gocache/consistenthash"
	pb "github.com/454270186/GoCache/gocache/gocachepb/gocachepb"
	"google.golang.org/protobuf/proto"
)

const (
	DefaultBasePath = "/_gocache/"
	HealthCheckPath = "/_check/"
	DefaultReplicas = 50

	HealthCheckInterval = 4 * time.Second
)

var _ PeerPicker = (*HTTPPool)(nil)

type HTTPPool struct {
	selfAddr    string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.HashRing // <data_key, peer_addr>
	peersAddrsSet map[string]bool
	httpGetters map[string]*httpGetter   // <peer_addr, the_Get()>
	httpPutters map[string]*httpPutter   // <peer_addr, the_Put()>
}

func NewHTTPPool(addr string) *HTTPPool {
	return &HTTPPool{
		selfAddr:    addr,
		basePath:    DefaultBasePath,
		httpGetters: make(map[string]*httpGetter),
		httpPutters: make(map[string]*httpPutter),
		peersAddrsSet: make(map[string]bool),
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s\n", p.selfAddr, fmt.Sprintf(format, v...))
}

// Get: 		/<BASEPATH>/<GroupName>/<Key>
// Put: 		/<BASEPATH>/<GroupName>/<Key>/<Val>
// HealthCheck: /<HealthCheck_PATH>/
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, DefaultBasePath) {
		if !strings.HasPrefix(r.URL.Path, HealthCheckPath) {
			http.Error(w, "[Error] bad base path", http.StatusNotFound)
			return
		} else {
			p.Log("[HealthCheck] %s", p.selfAddr)
			
			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var ping pb.Ping
			if err := proto.Unmarshal(bytes, &ping); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if ping.PingCode == 1 {
				// ping for checking health
				pong := pb.Pong{
					PongCode: 1,
				}
				pongResp, err := proto.Marshal(&pong)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				w.Header().Set("Content-Type", "octet-stream")
				w.Write(pongResp)
			}

			return
		}
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

		h.peersAddrsSet[addr] = true
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

func (h *HTTPPool) StartHealthCheck() {
	go func ()  {
		ticker := time.NewTicker(HealthCheckInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <- ticker.C:
				h.HealthCheck()
			}
		}
	}()
}

func (h *HTTPPool) HealthCheck() {
	for peerAddr := range h.peersAddrsSet {
		log.Printf("[CheckHealth] Check %v\n", peerAddr)

		ping := pb.Ping{
			PingCode: 1,
		}

		pingBody, err := proto.Marshal(&ping)
		if err != nil {
			log.Printf("[Error] %v\n", err.Error())
			return
		}

		pingBodyBuffer := bytes.NewBuffer(pingBody)
		resp, err := http.Post(peerAddr+HealthCheckPath, "", pingBodyBuffer)
		if err != nil {
			log.Printf("[Unhealthy] peer %v: error while send health check: %v\n", peerAddr, err.Error())
			h.peers.Remove(peerAddr)
			delete(h.peersAddrsSet, peerAddr)
			log.Printf("[Unhealthy] remove peer %v\n", peerAddr)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("[Unhealthy] peer %v\n: bad http resp code", peerAddr)
			h.peers.Remove(peerAddr)
			delete(h.peersAddrsSet, peerAddr)
			log.Printf("[Unhealthy] remove peer %v\n", peerAddr)
			continue
		}
		
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[Unhealthy] peer %v: error while read resp body: %v\n", peerAddr, err.Error())
			h.peers.Remove(peerAddr)
			delete(h.peersAddrsSet, peerAddr)
			log.Printf("[Unhealthy] remove peer %v\n", peerAddr)
			continue
		}

		var pong pb.Pong
		if err := proto.Unmarshal(bytes, &pong); err != nil {
			log.Printf("[Unhealthy] peer %v: error while unmarshal protobuf: %v\n", peerAddr, err.Error())
			h.peers.Remove(peerAddr)
			delete(h.peersAddrsSet, peerAddr)
			log.Printf("[Unhealthy] remove peer %v\n", peerAddr)
			return
		}

		if pong.PongCode != 1 {
			log.Printf("[Unhealthy] peer %v: invalid PongCode\n", peerAddr)
			h.peers.Remove(peerAddr)
			delete(h.peersAddrsSet, peerAddr)
			log.Printf("[Unhealthy] remove peer %v\n", peerAddr)
			continue
		}
	}
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
