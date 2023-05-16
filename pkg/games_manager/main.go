package main

import (
	"log"
	"net"
	"time"

	"github.com/rafaliusz/tictactoe/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	rpcTimeout = 5 * time.Second
)

func createGamesManagerServer() *gamesManagerServer {
	gms := &gamesManagerServer{}
	return gms
}

func startServer(gamesManagerServer *gamesManagerServer) {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	server.RegisterGamesManagerServer(grpcServer, gamesManagerServer)
	lis, err := net.Listen("tcp", "127.0.0.1:666")
	if err != nil {
		log.Fatalln(err.Error())
	}
	grpcServer.Serve(lis)
}

func main() {
	gamesManagerServer := createGamesManagerServer()
	startServer(gamesManagerServer)
}
