package node

import (
	"kvstore/server"
	"kvstore/storage"
	"log"
)

const (
	LEADER   = 1
	FOLLOWER = 2
	DELETE   = "delete"
	PUT      = "put"
)

type Node struct {
	id         int            // unique identifier for the node
	state      int            // current state of the node: LEADER or FOLLOWER LEADER = 1, FOLLOWER = 2
	leader     int            // id of the current leader if node is follower
	leaderAddr string         // address of the leader if node is follower
	nodes      map[int]string // ids of the followers if node is leader
	storage    *storage.MemoryStorage
	grpcServer *server.GRPCServer
}

func NewNode(id int, addr string, state int) *Node {
	n := &Node{
		id:      id,
		state:   state,
		leader:  -1, // no leader initially
		nodes:   make(map[int]string),
		storage: storage.NewMemoryStorage(),
	}

	n.grpcServer = server.NewServer(n)

	// Start server in goroutine so it doesn't block
	go func() {
		log.Printf("Node %d starting gRPC server on %s", id, addr)
		if err := n.grpcServer.StartGRPCServer(addr); err != nil {
			log.Printf("gRPC server error for node %d: %v", id, err)
		}
	}()

	return n
}

func (n *Node) GetID() int {
	return n.id
}

func (n *Node) GetState() int {
	return n.state
}

func (n *Node) IsLeader() bool {
	return n.state == LEADER
}

func (n *Node) SetLeader(leaderID int, addr string) {
	n.leader = leaderID
	n.leaderAddr = addr
	if leaderID == n.id {
		n.state = LEADER
	} else {
		n.state = FOLLOWER
	}
}

func (n *Node) AddFollower(followerID int, followerAddr string) {
	n.nodes[followerID] = followerAddr
}

func (n *Node) HandlePut(requesterID int, key, value string) error {
	if n.IsLeader() {
		err := n.storage.Put(key, value)
		if err != nil {
			return err
		}
		broadcastRequest(n.id, n.nodes, PUT, key, value)
	} else if requesterID != n.leader {
		forwardRequestToLeader(n.id, n.leaderAddr, PUT, key, value)
	} else {
		err := n.storage.Put(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) HandleGet(key string) (string, error) {
	value, err := n.storage.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (n *Node) HandleDelete(requesterID int, key string) error {
	if n.IsLeader() {
		_, err := n.storage.Delete(key)
		if err != nil {
			return err
		}
		broadcastRequest(n.id, n.nodes, DELETE, key, "")
	} else if requesterID != n.leader {
		forwardRequestToLeader(n.id, n.leaderAddr, DELETE, key, "")
	} else {
		_, err := n.storage.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
