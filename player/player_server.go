package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rafaliusz/tictactoe/pkg/game_server"
	"github.com/rafaliusz/tictactoe/pkg/logic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type playerServer struct {
	game_server.UnimplementedPlayerServer
	grpcServer    *grpc.Server
	port          int
	game          logic.TicTacToeGame
	playersSymbol logic.Symbol
	gameMutex     sync.RWMutex
	needMove      atomic.Bool
	move          chan [2]int
}

func (ps *playerServer) YourMove(ctx context.Context, empty *emptypb.Empty) (*game_server.YourMoveResult, error) {
	log.Println("YourMove")
	go ps.PlayerMove()
	return &game_server.YourMoveResult{}, nil
}

func (ps *playerServer) GameFinished(ctx context.Context, result *game_server.GameResult) (*game_server.GameFinishedResult, error) {
	log.Println("GameFinished")
	ps.gameMutex.Lock()
	defer ps.gameMutex.Unlock()
	if result.Result == game_server.GameResultEnum_Win {
		fmt.Println("You win. Congrats!")
	} else {
		fmt.Println("You lose. Better luck next time!")
	}

	return &game_server.GameFinishedResult{}, nil
}

func (server *playerServer) IsGameFinished() bool {
	server.gameMutex.Lock()
	defer server.gameMutex.Unlock()
	return server.game.GetGameState() != logic.InProgress
}

type MoveError struct{}

func (err MoveError) Error() string {
	return "Error occured during PlayerMove"
}

func (ps *playerServer) PlayerMove() {
	log.Println("PlayerMove")
	ps.gameMutex.Lock()
	defer ps.gameMutex.Unlock()
	ps.needMove.Store(true)
	move := <-ps.move

	conn, err := grpc.Dial("localhost:666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("PlayerMove dial error: %s", err.Error())
	}
	defer conn.Close()
	gamesManagerClient := game_server.NewGamesManagerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	md := metadata.New(map[string]string{"address": "localhost:" + strconv.Itoa(ps.port)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	defer cancel()
	res, err := gamesManagerClient.Move(ctx, &game_server.Position{Row: int32(move[0]), Column: int32(move[1])})
	if err != nil {
		log.Printf("Error returned by the server on Move: %s\n", err.Error())
		return
	}
	if res == nil {
		log.Fatalln("PlayerMove: nil Move result")
		return
	}

	if res.Result == game_server.MoveResultEnum_Ok {
		ps.game.Move(move[0], move[1], ps.playersSymbol)
		return
	}

	log.Printf("PlayerMove: move error: %s \n", err.Error())
}

func (ps *playerServer) UpdateGameState(ctx context.Context, position *game_server.Position) (*game_server.UpdateGameStateResult, error) {
	var opponentsSymbol logic.Symbol
	if ps.playersSymbol == logic.Circle {
		opponentsSymbol = logic.Cross
	} else {
		opponentsSymbol = logic.Circle
	}
	ps.game.Move(int(position.Row), int(position.Column), opponentsSymbol)
	return &game_server.UpdateGameStateResult{}, nil
}

func readPosition() (int, int) {
	for {
		fmt.Println("Enter your move:")
		var row, column int
		_, err := fmt.Scanf("%d %d", &row, &column)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		return row, column
	}
}
