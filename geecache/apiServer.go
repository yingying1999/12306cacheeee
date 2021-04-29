package geecache

import (
	"fmt"
	"log"
	"net/http"
)

type ApiServer struct {
	addr     string
	basePath string
	tm       *TicketManager
}

func NewApiServer(addr, basePath string, tm *TicketManager) *ApiServer {
	return &ApiServer{
		addr:     addr,
		basePath: basePath,
		tm:       tm,
	}
}

// ServeHTTP handle all http requests
func (p *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//输出用户请求
	p.Log("%s %s", r.Method, r.URL.Path)
	//获取参数
	key := r.URL.Query().Get("key")
	count := r.URL.Query().Get("count")
	startNo := r.URL.Query().Get("startNo")
	endNo := r.URL.Query().Get("endNo")
	//买票操作
	res, err := p.tm.BuyTickets(key, count, startNo, endNo)
	//响应
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// Log info with server name
func (p *ApiServer) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.addr, fmt.Sprintf(format, v...))
}
