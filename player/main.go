package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/rafaliusz/tictactoe/pkg/game_server"
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
	ps.address = listener.Addr().String()
	return ps
}

func tokenExists() bool {
	if _, err := os.Stat("token.txt"); err == nil {
		return true
	}
	return false
}

func getToken() (string, error) {
	dat, err := os.ReadFile("token.txt")
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func removeToken() {
	if tokenExists() {
		os.Remove("token.txt")
	}
}

func reconnect(ps *playerServer) bool {
	if !tokenExists() {
		return false
	}
	token, err := getToken()
	if err != nil {
		log.Printf("error reading token: %s", err.Error())
		return false
	}
	res, err := ps.reconnect(token)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !res {
		log.Printf("Couldn't reconnect")
		return false
	}
	return true
}

func join(ps *playerServer) {
	joined, err := ps.joinGame()
	if err != nil {
		log.Fatalf("Error joining: %s", err.Error())
	}
	if !joined {
		log.Fatalln("Error joining game")
	}
	log.Println("Joined the game")
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalln(err.Error())
	}
	ps := createPlayerServer(listener)
	go startServer(ps, listener)

	res := reconnect(ps)
	if !res {
		removeToken()
		join(ps)
	}

	for {
		if ps.needMove.Load() {
			fmt.Print(ps.game.ToString())
			row, column := readPosition()
			ps.needMove.Store(false)
			ps.move <- [2]int{row, column}
			continue
		}
		if ps.IsGameFinished() {
			removeToken()
			break
		}
		time.Sleep(time.Second * 2)
	}
	ps.grpcServer.GracefulStop()
}
