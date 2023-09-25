package gocache

import pb "github.com/454270186/GoCache/gocache/gocachepb/gocachepb"

type PeerPicker interface {
	PickPeer(key string) (peerGetter PeerGetter, peerPutter PeerPutter, ok bool)
}

type PeerGetter interface {
	Get(in *pb.GetRequest, out *pb.GetResponse) error
}

type PeerPutter interface {
	Put(in *pb.PutRequest, out *pb.PutResponse) error
}