package main

import (
	"fmt"
	"grpc_api/api"
	"grpc_api/database"
	"grpc_api/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("running grpc server on port 7777")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := api.Server{
		Database: database.NewDatabase("user"),
	}

	grpcserver := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcserver, &s)

	if err := grpcserver.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
