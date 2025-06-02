package main

import (
	"fmt"
	"kvstore/node"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Starting KV Store Node...")

	// Should me managed by a cluster manager
	n := node.NewNode(1, "localhost:50051", node.LEADER)
	n2 := node.NewNode(2, "localhost:50052", node.FOLLOWER)
	n3 := node.NewNode(3, "localhost:50053", node.FOLLOWER)
	n.AddFollower(2, "localhost:50052")
	n.AddFollower(3, "localhost:50053")
	n2.SetLeader(1)
	n3.SetLeader(1)

	log.Printf("Node %d created as %s", n.GetID(), getStateString(n.GetState()))
	log.Printf("Node %d created as %s", n2.GetID(), getStateString(n2.GetState()))
	log.Printf("Node %d created as %s", n3.GetID(), getStateString(n3.GetState()))

	time.Sleep(100 * time.Millisecond)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println(`
	Node is running on localhost:50051
	You can test it using grpcurl or a gRPC client
	
	Example grpcurl commands:
	# Put a key-value pair:
	grpcurl -plaintext -d '{\"key\":\"hello\",\"value\":\"world\"}' localhost:50051 KVStore/Put
	# Get a value:
	grpcurl -plaintext -d '{\"key\":\"hello\"}' localhost:50051 KVStore/Get
	# Delete a key:
	grpcurl -plaintext -d '{\"key\":\"hello\"}' localhost:50051 KVStore/Delete

	Press Ctrl+C to stop the node
	`)

	<-sigChan
	fmt.Println("\nReceived shutdown signal... stopping node")
	fmt.Println("Node stopped.")
}

func getStateString(state int) string {
	switch state {
	case node.LEADER:
		return "LEADER"
	case node.FOLLOWER:
		return "FOLLOWER"
	default:
		return "UNKNOWN"
	}
}
