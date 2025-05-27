package main

import (
	"fmt"
	"kvstore/server"
	"kvstore/storage"
)

func main() {
	fmt.Println("Hello, world!")

	storage := storage.NewMemoryStorage()
	fmt.Println("storage initialized")

	s := server.GRPCServer{
		Storage: storage,
	}

	err := s.StartGRPCServer("localhost:50051")
	if err != nil {
		panic(err)
	}
}
