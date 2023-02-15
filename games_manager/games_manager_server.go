package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/rafaliusz/tictactoe/pkg/game_server"
	"github.com/rafaliusz/tictactoe/pkg/logic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type gamesManagerServer struct {
	game_server.UnimplementedGamesManagerServer
	players      [2]PlayerInfo
	playersCount int32
	game         logic.TicTacToeGame
	gameMutex    sync.RWMutex
}

type PlayerInfo struct {
	address string
	token   string
	symbol  logic.Symbol
}

type JoinError string

func (err JoinError) Error() string {
	return string(err)
}

func (gs *gamesManagerServer) Join(ctx context.Context, in *empty.Empty) (*game_server.JoinResult, error) {
	gs.gameMutex.Lock()
	defer gs.gameMutex.Unlock()
	log.Println("Join called")
	if ctx == nil {
		return &game_server.JoinResult{Result: false, Info: "Internal error"}, JoinError("Nil context")
	}
	if gs.playersCount == 2 {
		log.Println("Join: Lobby full")
		return &game_server.JoinResult{Result: false, Info: "Lobby is full"}, nil
	}
	var symbol logic.Symbol
	if gs.playersCount == 0 {
		symbol = logic.Circle
	} else {
		symbol = logic.Cross
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &game_server.JoinResult{Result: false, Info: "Internal error"}, JoinError("Can't get metadata")
	}
	addressMD, ok := md["address"]
	if !ok {
		return &game_server.JoinResult{Result: false, Info: "Internal error"}, JoinError("Can't get address from metadata")
	}
	address := addressMD[0]
	token := uuid.New()

	log.Printf("Creating player: address %s\n", address)
	gs.players[gs.playersCount] = PlayerInfo{address: address, symbol: symbol, token: token.String()}
	gs.playersCount++
	if gs.playersCount == 2 {
		go gs.StartGame()
	}
	return &game_server.JoinResult{Result: true, Token: token.String(), Info: "Welcome to the lobby", Symbol: game_server.SymbolEnum(gs.players[gs.playersCount-1].symbol)}, nil
}

func (gs *gamesManagerServer) StartGame() {
	log.Println("StartGame")
	gs.YourMove(&gs.players[0])
}

func createPlayerClient(address string, timeout time.Duration) (game_server.PlayerClient, *grpc.ClientConn, *context.CancelFunc, *context.Context, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	playerClient := game_server.NewPlayerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return playerClient, conn, &cancel, &ctx, nil
}

func (gs *gamesManagerServer) UpdateGameState(player *PlayerInfo, position *game_server.Position) {
	log.Printf("UpdateGameState, player address: %s\n", player.address)
	gs.gameMutex.Lock()
	defer gs.gameMutex.Unlock()
	playerClient, connection, cancel, ctx, err := createPlayerClient(player.address, rpcTimeout)
	if err != nil {
		log.Fatalf("UpdateGameState: cannot create client: %s", err.Error())
		return
	}
	defer connection.Close()
	defer (*cancel)()
	_, err = (playerClient).UpdateGameState(*ctx, position)
	if err != nil {
		log.Printf("UpdateGameState error sending request: %s", err.Error())
	}
	if gs.game.GetGameState() == logic.InProgress {
		go gs.YourMove(player)
	}
}

func (gs *gamesManagerServer) YourMove(player *PlayerInfo) {
	log.Printf("YourMove, player address: %s\n", player.address)
	gs.gameMutex.Lock()
	defer gs.gameMutex.Unlock()
	playerClient, connection, cancel, ctx, err := createPlayerClient(player.address, rpcTimeout)
	if err != nil {
		log.Fatalf("YourMove: cannot create client: %s", err.Error())
		return
	}
	defer connection.Close()
	defer (*cancel)()
	_, err = (playerClient).YourMove(*ctx, &emptypb.Empty{})
	if err != nil {
		log.Printf("YourMove error sending request: %s", err.Error())
	}
}

type MoveError string

func (err MoveError) Error() string {
	return string(err)
}

type GameResult byte

const (
	Win GameResult = iota
	Loss
)

func (gs *gamesManagerServer) Move(ctx context.Context, position *game_server.Position) (*game_server.MoveResult, error) {
	fmt.Println("Move called")
	gs.gameMutex.Lock()
	defer gs.gameMutex.Unlock()
	if ctx == nil {
		return &game_server.MoveResult{Result: game_server.MoveResultEnum_Error}, MoveError("Nil context")
	}
	address, err := getFromMetadata(&ctx, "address")
	if err != nil {
		return &game_server.MoveResult{Result: game_server.MoveResultEnum_Error}, err
	}
	token, err := getFromMetadata(&ctx, "token")
	if err != nil {
		return &game_server.MoveResult{Result: game_server.MoveResultEnum_Error}, err
	}
	currentPlayer, otherPlayer := getPlayers(gs.players, token)
	if currentPlayer == nil || otherPlayer == nil {
		log.Printf("Invalid token provided: %s, address: %s", token, address)
		return &game_server.MoveResult{Result: game_server.MoveResultEnum_Error}, err
	}
	currentPlayer.address = address

	gameState, err := gs.game.Move(int(position.Column), int(position.Row), currentPlayer.symbol)
	if err != nil {
		go gs.YourMove(currentPlayer)
		return &game_server.MoveResult{Result: game_server.MoveResultEnum_Retry}, err
	}

	if gameState == logic.InProgress {
		go gs.UpdateGameState(otherPlayer, &game_server.Position{Row: position.Row, Column: position.Column})
	} else {
		go gs.FinishTheGame(position, struct {
			*PlayerInfo
			GameResult
		}{currentPlayer, Win}, struct {
			*PlayerInfo
			GameResult
		}{otherPlayer, Loss})
	}

	return &game_server.MoveResult{Result: game_server.MoveResultEnum_Ok}, nil
}

func getFromMetadata(ctx *context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(*ctx)
	if !ok {
		return "", MoveError("Can't get metadata")
	}
	value, ok := md[key]
	if !ok {
		return "", MoveError("Can't get from metadata: " + key)
	}
	return value[0], nil
}

func getPlayers(players [2]PlayerInfo, token string) (*PlayerInfo, *PlayerInfo) {
	var currentPlayer *PlayerInfo
	var otherPlayer *PlayerInfo

	if players[0].token == token {
		currentPlayer = &players[0]
		otherPlayer = &players[1]
	} else if players[1].token == token {
		currentPlayer = &players[1]
		otherPlayer = &players[0]
	}

	return currentPlayer, otherPlayer
}

func (gs *gamesManagerServer) FinishTheGame(position *game_server.Position, player1 struct {
	*PlayerInfo
	GameResult
}, player2 struct {
	*PlayerInfo
	GameResult
}) {
	gs.UpdateGameState(player2.PlayerInfo, &game_server.Position{Row: position.Row, Column: position.Column})
	gs.GameFinished(player1.PlayerInfo, player1.GameResult)
	gs.GameFinished(player2.PlayerInfo, player2.GameResult)
	gs.Reset()
}

func (gs *gamesManagerServer) GameFinished(player *PlayerInfo, gameResult GameResult) error {
	log.Println("GameFinished")
	gs.gameMutex.Lock()
	defer gs.gameMutex.Unlock()
	playerClient, connection, cancel, ctx, err := createPlayerClient(player.address, rpcTimeout)
	if err != nil {
		return err
	}
	defer connection.Close()
	defer (*cancel)()
	gameResultEnum := getGameResultEnum(gameResult)
	_, err = playerClient.GameFinished(*ctx, &game_server.GameResult{Result: gameResultEnum})
	return err
}

func getGameResultEnum(gameResult GameResult) game_server.GameResultEnum {
	if gameResult == Win {
		return game_server.GameResultEnum_Win
	}
	return game_server.GameResultEnum_Loss
}

func (gs *gamesManagerServer) Reset() {
	gs.game = logic.TicTacToeGame{}
	gs.playersCount = 0
	gs.players = [2]PlayerInfo{}
}

func (gs *gamesManagerServer) Reconnect(ctx context.Context, reconnectData *game_server.ReconnectData) (*game_server.ReconnectResult, error) {
	log.Println("Reconnect")
	gs.gameMutex.Lock()
	defer gs.gameMutex.Unlock()
	if gs.playersCount < 2 {
		log.Printf("Reconnect: expired token, there's no active game: %s\n", reconnectData.Token)
		return &game_server.ReconnectResult{Result: false, Text: "expired token, there's no active game"}, nil
	}
	player := getPlayer(&gs.players, reconnectData.Token)
	if player == nil {
		log.Printf("Reconnect: invalid token: %s\n", reconnectData.Token)
		return &game_server.ReconnectResult{Result: false, Text: "invalid token"}, nil
	}
	address, err := getFromMetadata(&ctx, "address")
	if err != nil {
		log.Fatalf("Reconnect: error getting address from metadata: %s\n", err.Error())
		return &game_server.ReconnectResult{Result: false, Text: "error retrieving players's address"}, err
	}
	player.address = address
	if player.symbol == gs.game.GetNextMove() {
		go gs.YourMove(player)
	}
	gameState := gs.game.ToByteArray()
	return &game_server.ReconnectResult{Result: true, Text: "you've reconnected successfully", Board: gameState[:], Symbol: game_server.SymbolEnum(player.symbol)}, nil
}

func getPlayer(players *[2]PlayerInfo, token string) *PlayerInfo {
	if players[0].token == token {
		return &players[0]
	}
	if players[1].token == token {
		return &players[1]
	}
	return nil
}
