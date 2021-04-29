package geecache

// PeerPicker可以通过key得到对应的PeerUpdater
type PeerPicker interface {
	PickPeer(key string) (peer PeerUpdater, ok bool)
}

type PeerUpdater interface {
	Update(key, count, startNo, endNo string) ([]byte, error)
}
