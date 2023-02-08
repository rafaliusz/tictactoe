package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/rafaliusz/tictactoe/pkg/game_server"
	"github.com/rafaliusz/tictactoe/pkg/logic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	serverAddress = "localhost:666"
)

func createServer() *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	return grpcServer
}

func startServer(server *playerServer) {
	server.grpcServer = createServer()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", server.port))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	game_server.RegisterPlayerServer(server.grpcServer, server)
	server.grpcServer.Serve(lis)
}

func joinGame(port int) (bool, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer conn.Close()
	gamesManagerClient := game_server.NewGamesManagerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	md := metadata.New(map[string]string{"address": "localhost:" + strconv.Itoa(port)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	defer cancel()
	res, err := gamesManagerClient.Join(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}
	return res.Result, nil
}

func createPlayerServer(port int) *playerServer {
	var ps playerServer
	ps.port = port
	ps.move = make(chan [2]int, 1)
	ps.playersSymbol = logic.Circle
	return &ps
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Need port as first argument")
		os.Exit(1)
	}
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid port")
		os.Exit(1)
	}
	joined, err := joinGame(port)
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
	ps := createPlayerServer(port)
	go startServer(ps)
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
}
