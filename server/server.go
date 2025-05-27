package server

import (
	"context"
	"fmt"
	"kvstore/storage"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	Storage *storage.MemoryStorage
	UnimplementedKVStoreServer
}

func NewServer(storage *storage.MemoryStorage) *GRPCServer {
	return &GRPCServer{Storage: storage}
}

func (s *GRPCServer) Put(ctx context.Context, req *PutRequest) (*PutResponse, error) {
	err := s.Storage.Put(req.Key, req.Value)
	return &PutResponse{Success: err == nil}, err
}

func (s *GRPCServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	value, err := s.Storage.Get(req.Key)
	return &GetResponse{Value: value}, err
}

func (s *GRPCServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	bool, err := s.Storage.Delete(req.Key)
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
	RegisterKVStoreServer(grpcServer, &GRPCServer{Storage: s.Storage})
	return grpcServer.Serve(listener)
}
