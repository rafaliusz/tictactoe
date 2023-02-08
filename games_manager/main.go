package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/rafaliusz/tictactoe/pkg/game_server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 666))
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
