package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/rafaliusz/tictactoe/pkg/game_server"
	"github.com/rafaliusz/tictactoe/pkg/logic"
	"google.golang.org/grpc"
)

const (
	serverAddress = "localhost:666"
	joinTimeout   = 5 * time.Second
	moveTimeout   = 5 * time.Second
)

func createServer() *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	return grpcServer
}

func startServer(server *playerServer, listener net.Listener) {
	server.address = listener.Addr().String()
	game_server.RegisterPlayerServer(server.grpcServer, server)
	log.Printf("starting server: %s", listener.Addr().String())
	err := server.grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("serve error: %s", err.Error())
	}
}

func createPlayerServer(listener net.Listener) *playerServer {
	ps := &playerServer{}
	ps.grpcServer = createServer()
	ps.move = make(chan [2]int, 1)
	ps.playersSymbol = logic.Circle
	ps.address = listener.Addr().String()
	return ps
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	ps := createPlayerServer(listener)
	go startServer(ps, listener)
	joined, err := ps.joinGame()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	if !joined {
		log.Println("Could not join the game")
		os.Exit(-1)
	} else {
		log.Println("Joined the game")
	}
	log.Println("server started")

	for {
		if ps.needMove.Load() {
			fmt.Print(ps.game.ToString())
			row, column := readPosition()
			fmt.Println("main: got move")
			ps.needMove.Store(false)
			fmt.Println("main: need move set")
			ps.move <- [2]int{row, column}
			fmt.Println("main: move sent")
			continue
		}
		if ps.IsGameFinished() {
			break
		}
		time.Sleep(time.Second * 2)
	}
	ps.grpcServer.GracefulStop()
}
