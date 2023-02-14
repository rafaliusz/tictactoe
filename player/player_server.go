package main

import (
	"context"
	"fmt"
	"log"
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
	address       string
	game          logic.TicTacToeGame
	playersSymbol logic.Symbol
	gameMutex     sync.RWMutex
	needMove      atomic.Bool
	move          chan [2]int
	token         string
}

func createGamesManagerClient(address string, token string, timeout time.Duration) (game_server.GamesManagerClient, *grpc.ClientConn, *context.CancelFunc, *context.Context, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	gamesManagerClient := game_server.NewGamesManagerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	metadataMap := map[string]string{"address": address}
	if token != "" {
		metadataMap["token"] = token
	}
	md := metadata.New(metadataMap)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return gamesManagerClient, conn, &cancel, &ctx, nil
}

func (ps *playerServer) joinGame() (bool, error) {
	gamesManagerClient, connection, cancel, ctx, err := createGamesManagerClient(ps.address, ps.token, joinTimeout)
	if err != nil {
		return false, err
	}
	defer connection.Close()
	defer (*cancel)()
	res, err := gamesManagerClient.Join(*ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}
	if res.Result {
		ps.token = res.Token
	}
	return res.Result, nil
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

	gamesManagerClient, connection, cancel, ctx, err := createGamesManagerClient(ps.address, ps.token, moveTimeout)
	if err != nil {
		log.Fatalf("PlayerMove: error creating client: %s", err.Error())
		return
	}
	defer connection.Close()
	defer (*cancel)()
	res, err := gamesManagerClient.Move(*ctx, &game_server.Position{Row: int32(move[0]), Column: int32(move[1])})
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

	log.Println("PlayerMove: server returned error")
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
