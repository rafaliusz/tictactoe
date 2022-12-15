package main

import "fmt"

func main() {
	var tictactoe TicTacToeGame
	tictactoe.move(1, 1, Cross)
	tictactoe.move(2, 2, Circle)
	tictactoe.move(2, 0, Cross)

	_, err := tictactoe.move(2, 0, Cross)
	if err != nil {
		fmt.Println(err)
	}
	_, err = tictactoe.move(3, 0, Circle)
	if err != nil {
		fmt.Println(err)
	}
	gameState, err := tictactoe.move(0, 0, Circle)
	fmt.Println(gameState)

	gameState, err = tictactoe.move(0, 2, Cross)
	fmt.Println(gameState)

	fmt.Println(tictactoe.board.toString())

	_, err = tictactoe.move(0, 0, Circle)
	if err != nil {
		fmt.Println(err)
	}
}
