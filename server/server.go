package server

import (
	"context"
	"kvstore/iface"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	node iface.NodeAPI
	UnimplementedKVStoreServer
}

func NewServer(n iface.NodeAPI) *GRPCServer {
	return &GRPCServer{node: n}
}

func (s *GRPCServer) Put(ctx context.Context, req *PutRequest) (*PutResponse, error) {
	var requesterID int
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["requesterid"]; ok && len(val) > 0 {
			id, err := strconv.Atoi(val[0])
			if err == nil {
				requesterID = id
			}
		}
	}

	err := s.node.HandlePut(requesterID, req.Key, req.Value)
	return &PutResponse{Success: err == nil}, err
}

func (s *GRPCServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	value, err := s.node.HandleGet(req.Key)
	return &GetResponse{Value: value}, err
}

func (s *GRPCServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	var requesterID int
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["requesterid"]; ok && len(val) > 0 {
			id, err := strconv.Atoi(val[0])
			if err == nil {
				requesterID = id
			}
		}
	}
	err := s.node.HandleDelete(requesterID, req.Key)

	return &DeleteResponse{Success: err == nil}, err
}

func (s *GRPCServer) StartGRPCServer(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	RegisterKVStoreServer(grpcServer, &GRPCServer{node: s.node})
	return grpcServer.Serve(listener)
}
