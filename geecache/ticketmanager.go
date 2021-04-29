package geecache

import (
	"encoding/json"
	"errors"
	"geecache/concurrentmap"
	"log"
	"strconv"
)

//每个缓存服务器一个ticketmanager
// A Group is a cache namespace and associated data loaded spread over
type TicketManager struct {
	concurrentMap concurrentmap.ConcurrentMap
	peerPicker    PeerPicker
}

func NewTM() *TicketManager {
	tm := &TicketManager{
		concurrentMap: concurrentmap.NewConcurrentMap(),
	}
	return tm
}

func (tm *TicketManager) SetPeerPicker(peerPicker PeerPicker) {
	if tm.peerPicker != nil {
		panic("RegisterPeerPicker called more than once")
	}
	tm.peerPicker = peerPicker
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
		return tm.RemoteUpdate(key, count, startNo, endNo)
	}
}

func (tm *TicketManager) RemoteUpdate(key, count, startNo, endNo string) ([]byte, error) {
	if tm.peerPicker != nil {
		if peerUpdater, ok := tm.peerPicker.PickPeer(key); ok {
			responseObj, err := peerUpdater.Update(key, count, startNo, endNo)
			if err != nil {
				return nil, err
			}
			return responseObj, nil
		}
	}
	return nil, errors.New("不存在该key")
}
