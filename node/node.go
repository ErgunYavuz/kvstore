package node

import (
	"context"
	"kvstore/server"
	"kvstore/storage"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	LEADER   = 1
	FOLLOWER = 2
)

type Node struct {
	id         int            // unique identifier for the node
	state      int            // current state of the node: LEADER or FOLLOWER LEADER = 1, FOLLOWER = 2
	leader     int            // id of the current leader if node is follower
	followers  map[int]string // ids of the followers if node is leader
	storage    *storage.MemoryStorage
	grpcServer *server.GRPCServer
}

func NewNode(id int, addr string, state int) *Node {
	n := &Node{
		id:        id,
		state:     state,
		leader:    -1, // no leader initially
		followers: make(map[int]string),
		storage:   storage.NewMemoryStorage(),
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

func (n *Node) SetLeader(leaderID int) {
	n.leader = leaderID
	if leaderID == n.id {
		n.state = LEADER
	} else {
		n.state = FOLLOWER
	}
}

func (n *Node) AddFollower(followerID int, followerAddr string) {
	n.followers[followerID] = followerAddr
}

func (n *Node) HandlePut(key, value string) error {
	err := n.storage.Put(key, value)
	if err != nil {
		return err
	}

	if n.IsLeader() {
		for followerID, followerAddr := range n.followers {
			go func(id int, addr string) {
				log.Printf("Node %d replicating Put to follower %d at %s", n.id, id, addr)
				conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("Failed to connect to follower %d: %v", id, err)
					return
				}
				defer conn.Close()

				client := server.NewKVStoreClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				_, err = client.Put(ctx, &server.PutRequest{Key: key, Value: value})
				if err != nil {
					log.Printf("Replication to follower %d failed: %v", id, err)
				} else {
					log.Printf("Successfully replicated to follower %d", id)
				}
			}(followerID, followerAddr)
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

func (n *Node) Delete(key string) (bool, error) {
	success, err := n.storage.Delete(key)
	if err != nil {
		return false, err
	}

	if n.IsLeader() {
		for followerID, followerAddr := range n.followers {
			go func(id int, addr string) {
				log.Printf("Node %d replicating Put to follower %d at %s", n.id, id, addr)
				conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("Failed to connect to follower %d: %v", id, err)
					return
				}
				defer conn.Close()

				client := server.NewKVStoreClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				_, err = client.Delete(ctx, &server.DeleteRequest{Key: key})
				if err != nil {
					log.Printf("Replication to follower %d failed: %v", id, err)
				} else {
					log.Printf("Successfully replicated to follower %d", id)
				}
			}(followerID, followerAddr)
		}
	}
	return success, nil
}
