package geecache

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
)

//每个缓存服务器一个ticketmanager
// A Group is a cache namespace and associated data loaded spread over
type TicketManager struct {
	concurrentMap ConcurrentMap
}

var instance *TicketManager

func GetTicketManager() *TicketManager {
	if instance == nil {
		instance = NewTM() // <--- NOT THREAD SAFE
	}
	return instance
}
func NewTM() *TicketManager {
	tm := &TicketManager{
		concurrentMap: NewConcurrentMap(),
	}
	return tm
}
func (tm *TicketManager) GetMap() ConcurrentMap {
	return tm.concurrentMap
}

func (tm *TicketManager) BuyTickets(key, count, startNo, endNo string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("key不允许为空")
	}

	if _, ok := tm.concurrentMap.Get(key); ok {
		log.Println("[GeeCache] hit")
		countInt, _ := strconv.Atoi(count)
		startNoInt, _ := strconv.Atoi(startNo)
		endNoInt, _ := strconv.Atoi(endNo)
		res, err := tm.concurrentMap.Update(key, int32(countInt), int32(startNoInt), int32(endNoInt))
		resBytes, _ := json.Marshal(res)
		return resBytes, err
	} else {
		log.Println("[GeeCache] miss")
		return nil, errors.New("miss")
	}
}
