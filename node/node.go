package node

import (
	"kvstore/server"
	"kvstore/storage"
)

const (
	LEADER   = 1
	FOLLOWER = 2
)

type Node struct {
	id         int
	state      int
	leader     int   //id of the current leader if node is follower
	followers  []int //ids of the followers if node is leader
	grpcServer server.GRPCServer
}

func NewNode(id int, addr string, state int) *Node {
	n := &Node{
		id:        id,
		state:     state,
		leader:    -1, // no leader initially
		followers: []int{},
	}

	n.grpcServer = *server.NewServer(storage.NewMemoryStorage())
	n.grpcServer.StartGRPCServer(addr)
	return n
}

// upon receive put
//    if state is leader
// 	      write put
// 	      broadcast to followers //replication
