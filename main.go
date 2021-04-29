package main

import (
	"flag"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func startCacheServer(addr, basePath string, tm *geecache.TicketManager) {
	log.Println("cacheServer is running at", addr)
	cacheServer := geecache.NewCacheServer(addr, "addr", tm)
	log.Fatal(http.ListenAndServe(addr[7:], cacheServer))
}

func startAPIServer(apiAddr, basePath string, tm *geecache.TicketManager) {
	log.Println("APIServer is running at", apiAddr)
	apiServer := geecache.NewApiServer(apiAddr, "addr", tm)
	log.Fatal(http.ListenAndServe(apiAddr[7:], apiServer))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	basePath := "buyticket"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	//初始化peerPicker
	peerPicker := geecache.NewHTTPPeerPicker(addrMap[port], "/_geecache/")
	peerPicker.Set(addrs...)
	//初始化ticketManager
	tm := geecache.NewTM()
	tm.SetPeerPicker(peerPicker)
	//开启api服务
	if api {
		go startAPIServer(apiAddr, basePath, tm)
	}
	//开启缓存服务
	startCacheServer(addrMap[port], basePath, tm)
}
