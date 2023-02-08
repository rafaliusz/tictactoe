package logic

import "fmt"

type TicTacToeGame struct {
	Board
	lastMove  Symbol
	gameState GameState
}

type GameState byte

const (
	InProgress GameState = iota
	CrossWins
	CircleWins
)

func (game *TicTacToeGame) GetNextMove() Symbol {
	if game.lastMove == Cross || game.lastMove == None {
		return Circle
	}
	return Cross
}

type InvalidMoveError string

func (err InvalidMoveError) Error() string {
	return string(err)
}

func (game *TicTacToeGame) Move(row int, column int, symbol Symbol) (GameState, error) {
	if game.gameState != InProgress {
		return game.gameState, InvalidMoveError("Cannot move when the game is finished")
	}
	if symbol == None {
		return game.gameState, InvalidMoveError("Cannot use None as a move symbol")
	}
	if symbol == game.lastMove {
		return game.gameState, InvalidMoveError(fmt.Sprintf("Cannot use %s twice in a row", symbol.String()))
	}

	err := move(&game.Board, row, column, symbol)
	if err != nil {
		return game.gameState, err
	}
	game.lastMove = symbol
	game.gameState = getGameState(&game.Board)
	return game.gameState, nil
}

func (game *TicTacToeGame) GetGameState() GameState {
	return game.gameState
}

func move(board *Board, row int, column int, symbol Symbol) error {
	if row > 2 || column > 2 {
		return InvalidMoveError(fmt.Sprintf("Invalid index, %d:%d is not within the boundaries of the board", row, column))
	}
	if board[row][column] != None {
		return InvalidMoveError(fmt.Sprintf("Invalid move, %d:%d is being used by %s", row, column, (&symbol).Name()))
	}
	board[row][column] = symbol
	return nil
}

func getGameState(board *Board) GameState {

	if board[1][1] != None {
		gameState := checkColumn(board, 1)
		if gameState != InProgress {
			return gameState
		}
		gameState = checkRow(board, 1)
		if gameState != InProgress {
			return gameState
		}
		gameState = checkDiagonals(board)
		if gameState != InProgress {
			return gameState
		}
	}

	if board[1][0] != None {
		gameState := checkColumn(board, 0)
		if gameState != InProgress {
			return gameState
		}
	}
	if board[1][2] != None {
		gameState := checkColumn(board, 2)
		if gameState != InProgress {
			return gameState
		}
	}
	if board[0][1] != None {
		gameState := checkRow(board, 0)
		if gameState != InProgress {
			return gameState
		}
	}
	if board[2][1] != None {
		gameState := checkRow(board, 2)
		if gameState != InProgress {
			return gameState
		}
	}
	return InProgress
}

func symbolToGameState(symbol Symbol) GameState {
	switch symbol {
	case None:
		return InProgress
	case Cross:
		return CrossWins
	case Circle:
		return CircleWins
	default:
		return InProgress
	}
}

func checkRow(board *Board, num int) GameState {
	symbol := board[num][0]
	for i := 1; i < len(board[num]); i++ {
		if symbol != board[num][i] {
			return InProgress
		}
	}
	return symbolToGameState(symbol)
}

func checkColumn(board *Board, num int) GameState {
	symbol := board[0][num]
	for i := 1; i < len(board); i++ {
		if symbol != board[i][num] {
			return InProgress
		}
	}
	return symbolToGameState(symbol)
}

func checkDiagonals(board *Board) GameState {

	symbol := board[0][0]
	equal := true
	for i := 1; i < len(board); i++ {
		if symbol != board[i][i] {
			equal = false
			break
		}
	}
	if equal {
		return symbolToGameState(symbol)
	}
	symbol = board[len(board)-1][0]
	for i := 1; i < len(board); i++ {
		if symbol != board[len(board)-1-i][i] {
			return InProgress
		}
	}
	return symbolToGameState(symbol)
}
