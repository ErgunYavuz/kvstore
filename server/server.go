package server

import (
	"context"
	"fmt"
	"kvstore/iface"
	"net"

	"google.golang.org/grpc"
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
	err := s.node.HandlePut(req.Key, req.Value)
	return &PutResponse{Success: err == nil}, err
}

func (s *GRPCServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	value, err := s.node.HandleGet(req.Key)
	return &GetResponse{Value: value}, err
}

func (s *GRPCServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	bool, err := s.node.Delete(req.Key)
	return &DeleteResponse{Success: bool}, err
}

func (s *GRPCServer) StartGRPCServer(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Errorf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	RegisterKVStoreServer(grpcServer, &GRPCServer{node: s.node})
	return grpcServer.Serve(listener)
}
