package gocache

type PeerPicker interface {
	PickPeer(key string) (peerGetter PeerGetter, peerPutter PeerPutter, ok bool)
}

type PeerGetter interface {
	Get(group, key string) (string, error)
}

type PeerPutter interface {
	Put(group, key, val string) error
}