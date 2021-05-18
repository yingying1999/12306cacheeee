package main

import (
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
	// 新建服务器的handler
	cacheServer := geecache.NewCacheServer(addr, "addr", tm)
	// 启动服务器
	log.Fatal(http.ListenAndServe(addr[7:], cacheServer))
}

func main() {
	addr := "http://localhost:8001"
	basePath := "buyticket"

	//初始化ticketManager
	tm := geecache.GetTicketManager()
	InitCache()
	//开启缓存服务
	startCacheServer(addr, basePath, tm)
}

// 从数据库中读取车次信息，并初始化对应的tickets
func InitCache() {
	myMap := geecache.GetTicketManager().GetMap()
	myMap.Init("C2201-2020_02_06-BusinessSeat", 10, 32)
}
