package main

import "testing"

func TestMoveOutsideBoardsBoundaries(t *testing.T) {
	var game TicTacToeGame
	boardSize := len(game.board)
	_, err := game.move(boardSize, 0, Circle)
	if err == nil {
		t.Error("Move outside boards boundaries should result in error")
	}
}

func TestMoveSameSymbolTwiceInARow(t *testing.T) {
	var game TicTacToeGame
	game.move(0, 0, Circle)
	_, err := game.move(0, 1, Circle)
	if err == nil {
		t.Error("Using same symbol twice in a row should result in error")
	}
}

func TestMoveSamePositionTwice(t *testing.T) {
	var game TicTacToeGame
	game.move(0, 0, Circle)
	_, err := game.move(0, 0, Cross)
	if err == nil {
		t.Error("Moving to an already taken position should result in error")
	}
}

func TestMoveWithNoneSymbol(t *testing.T) {
	var game TicTacToeGame
	_, err := game.move(0, 0, None)
	if err == nil {
		t.Error("Moving with None symbol should result in error")
	}
}

func TestMoveOnAFinishedBoard(t *testing.T) {
	var game TicTacToeGame
	game.move(0, 0, Circle)
	game.move(1, 0, Cross)
	game.move(0, 1, Circle)
	game.move(1, 1, Cross)
	game.move(0, 2, Circle)
	gameState, err := game.move(0, 2, Cross)
	if gameState != CircleWins {
		t.Error("GameState should be CircleWins")
	}
	if err == nil {
		t.Error("Moving on a finished board should result in error")
	}
}

func finishGameRow(game *TicTacToeGame, symbol Symbol, num int) {
	var anotherSymbol Symbol
	if symbol == Circle {
		anotherSymbol = Cross
	} else {
		anotherSymbol = Circle
	}
	anotherRow := (num + 1) % len(game.board)
	for i := 0; i < len(game.board[0]); i++ {
		game.move(num, i, symbol)
		game.move(anotherRow, i, anotherSymbol)
	}
}

func TestMoveFinishInRows(t *testing.T) {
	for i := 0; i < 3; i++ {
		game := TicTacToeGame{}
		finishGameRow(&game, Circle, i)
		if game.gameState != CircleWins {
			t.Errorf("Circle should win in row %d", i)
		}
		game = TicTacToeGame{}
		finishGameRow(&game, Cross, i)
		if game.gameState != CrossWins {
			t.Errorf("Cross should win in row %d", i)
		}
	}
}

func finishGameColumn(game *TicTacToeGame, symbol Symbol, num int) {
	var anotherSymbol Symbol
	if symbol == Circle {
		anotherSymbol = Cross
	} else {
		anotherSymbol = Circle
	}
	anotherColumn := (num + 1) % len(game.board)
	for i := 0; i < len(game.board[0]); i++ {
		game.move(i, num, symbol)
		game.move(i, anotherColumn, anotherSymbol)
	}
}

func TestMoveFinishInColumns(t *testing.T) {
	for i := 0; i < 3; i++ {
		game := TicTacToeGame{}
		finishGameRow(&game, Circle, i)
		if game.gameState != CircleWins {
			t.Errorf("Circle should win in column %d", i)
		}
		game = TicTacToeGame{}
		finishGameRow(&game, Cross, i)
		if game.gameState != CrossWins {
			t.Errorf("Cross should win in column %d", i)
		}
	}
}
