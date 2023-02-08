package main

import (
	"fmt"

	"github.com/rafaliusz/tictactoe/pkg/logic"
)

func main() {
	var tictactoe logic.TicTacToeGame
	tictactoe.Move(1, 1, logic.Cross)
	tictactoe.Move(2, 2, logic.Circle)
	tictactoe.Move(2, 0, logic.Cross)

	_, err := tictactoe.Move(2, 0, logic.Cross)
	if err != nil {
		fmt.Println(err)
	}
	_, err = tictactoe.Move(3, 0, logic.Circle)
	if err != nil {
		fmt.Println(err)
	}
	gameState, err := tictactoe.Move(0, 0, logic.Circle)
	fmt.Println(gameState)

	gameState, err = tictactoe.Move(0, 2, logic.Cross)
	fmt.Println(gameState)

	fmt.Println(tictactoe.Board.ToString())

	_, err = tictactoe.Move(0, 0, logic.Circle)
	if err != nil {
		fmt.Println(err)
	}
}
