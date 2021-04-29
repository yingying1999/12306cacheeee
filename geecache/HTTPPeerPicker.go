package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

//httppool实现PeerPicker，可以通过key得到对应的PeerUpdater
//httppool还实现缓存服务器
const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPeerPicker struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	addr            string
	basePath        string
	mu              sync.Mutex // guards peers and httpGetters
	peersMap        *consistenthash.Map
	peerUpdatersMap map[string]*HTTPPeerUpdater // keyed by e.g. "http://10.0.0.2:8008"
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPeerPicker(addr, basePath string) *HTTPPeerPicker {
	return &HTTPPeerPicker{
		addr:     addr,
		basePath: basePath,
	}
}

// Set updates the pool's list of peers.
func (p *HTTPPeerPicker) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peersMap = consistenthash.New(defaultReplicas, nil)
	p.peersMap.Add(peers...)
	p.peerUpdatersMap = make(map[string]*HTTPPeerUpdater, len(peers))
	for _, peer := range peers {
		p.peerUpdatersMap[peer] = &HTTPPeerUpdater{baseURL: peer + p.basePath}
	}
}

// PickPeer picks a peer according to key
func (p *HTTPPeerPicker) PickPeer(key string) (PeerUpdater, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peersMap.Get(key); peer != "" && peer != p.addr {
		p.Log("Pick peer %s", peer)
		return p.peerUpdatersMap[peer], true
	}
	return nil, false
}

// Log info with server name
func (p *HTTPPeerPicker) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.addr, fmt.Sprintf(format, v...))
}

var _ PeerPicker = (*HTTPPeerPicker)(nil)

type HTTPPeerUpdater struct {
	baseURL string
}

func (h *HTTPPeerUpdater) Update(key, count, startNo, endNo string) ([]byte, error) {
	fmt.Println("向远程节点update")
	u := fmt.Sprintf(
		"%v?key=%v&count=%v&startNo=%v&endNo=%v/",
		h.baseURL, key, count, startNo, endNo,
	)
	fmt.Println("u", u)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ PeerUpdater = (*HTTPPeerUpdater)(nil)
