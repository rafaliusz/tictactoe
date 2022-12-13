package main

import "fmt"

func main() {
	var board Board
	board.move(1, 1, Cross)
	board.move(2, 2, Circle)
	board.move(2, 0, Cross)

	err := board.move(2, 0, Cross)
	if err != nil {
		fmt.Println(err)
	}
	err = board.move(3, 0, Cross)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(board.toString())
}
