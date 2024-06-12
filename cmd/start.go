package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"sajoniks.github.io/movieApi/internal/app"
	"sajoniks.github.io/movieApi/pkg/proto-gen"
)

func main() {
	api, err := app.NewMovieApiServer()
	if err != nil {
		log.Fatalf("failed to create API: %v", err)
	}
	lis, err := net.Listen("tcp", "localhost:55005")
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	opts := []grpc.ServerOption{}
	serv := grpc.NewServer(opts...)
	proto_gen.RegisterMovieApiServer(serv, api)
	serv.Serve(lis)
}
