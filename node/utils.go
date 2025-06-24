package node

import (
	"context"
	"kvstore/server"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func broadcastRequest(myid int, nodes map[int]string, requestType string, key string, value string) {
	switch requestType {
	case "put":
		for followerID, followerAddr := range nodes {
			if followerID == myid {
				log.Printf("Node %d is leader, skipping replication to itself", myid)
				continue
			}
			go func(id int, addr string) {
				log.Printf("Node %d replicating Put to follower %d at %s", myid, id, addr)
				conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("Failed to connect to follower %d: %v", id, err)
					return
				}
				defer conn.Close()

				client := server.NewKVStoreClient(conn)
				md := metadata.Pairs("requesterID", strconv.Itoa(myid))
				ctx := metadata.NewOutgoingContext(context.TODO(), md)
				ctx, cancel := context.WithTimeout(ctx, time.Second)
				defer cancel()

				_, err = client.Put(ctx, &server.PutRequest{Key: key, Value: value})
				if err != nil {
					log.Printf("Replication to follower %d failed: %v", id, err)
				} else {
					log.Printf("Successfully replicated to follower %d", id)
				}
			}(followerID, followerAddr)
		}
	case "delete":
		for followerID, followerAddr := range nodes {
			if followerID == myid {
				log.Printf("Node %d is leader, skipping replication to itself", myid)
				continue
			}
			go func(id int, addr string) {
				log.Printf("Node %d replicating delete to follower %d at %s", myid, id, addr)
				conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("Failed to connect to follower %d: %v", id, err)
					return
				}
				defer conn.Close()

				client := server.NewKVStoreClient(conn)
				md := metadata.Pairs("requesterID", strconv.Itoa(myid))
				ctx := metadata.NewOutgoingContext(context.TODO(), md)
				ctx, cancel := context.WithTimeout(ctx, time.Second)
				defer cancel()

				_, err = client.Delete(ctx, &server.DeleteRequest{Key: key})
				if err != nil {
					log.Printf("Replication to follower %d failed: %v", id, err)
				} else {
					log.Printf("Successfully replicated to follower %d", id)
				}
			}(followerID, followerAddr)
		}
	default:
		log.Printf("Unknown request type: %s", requestType)
	}
}

func forwardRequestToLeader(myid int, leaderAddr string, requestType string, key string, value string) {
	switch requestType {
	case "put":
		conn, err := grpc.NewClient(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect to leader %s: %v", leaderAddr, err)
			return
		}
		defer conn.Close()

		client := server.NewKVStoreClient(conn)
		md := metadata.Pairs("requesterID", strconv.Itoa(myid))
		ctx := metadata.NewOutgoingContext(context.TODO(), md)
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		_, err = client.Put(ctx, &server.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Printf("Forwarding Put to leader %s failed: %v", leaderAddr, err)
		} else {
			log.Printf("Successfully forwarded Put to leader %s", leaderAddr)
		}

	case "delete":
		conn, err := grpc.NewClient(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect to leader %s: %v", leaderAddr, err)
			return
		}
		defer conn.Close()

		client := server.NewKVStoreClient(conn)
		md := metadata.Pairs("requesterID", strconv.Itoa(myid))
		ctx := metadata.NewOutgoingContext(context.TODO(), md)
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		_, err = client.Delete(ctx, &server.DeleteRequest{Key: key})
		if err != nil {
			log.Printf("Forwarding delete to leader %s failed: %v", leaderAddr, err)
		} else {
			log.Printf("Successfully forwarded delete to leader %s", leaderAddr)
		}
	default:
		log.Printf("Unknown request type: %s", requestType)
	}
}
