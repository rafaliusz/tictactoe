package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/rafaliusz/tictactoe/pkg/game_server"
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
	game_server.RegisterGamesManagerServer(grpcServer, gamesManagerServer)
	lis, err := net.Listen("tcp", "127.0.0.1:666")
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(-1)
	}
	grpcServer.Serve(lis)
}

func main() {
	gamesManagerServer := createGamesManagerServer()
	startServer(gamesManagerServer)
}
