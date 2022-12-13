package main

import "fmt"

type Symbol int64

const (
	None Symbol = iota
	Circle
	Cross
)

type Board [3][3]Symbol
type InvalidMoveError string

func (symbol *Symbol) Name() string {
	switch *symbol {
	case None:
		return "None"
	case Circle:
		return "Circle"
	case Cross:
		return "Cross"
	default:
		return "Unknown"
	}
}

func (symbol *Symbol) String() string {
	switch *symbol {
	case None:
		return " "
	case Circle:
		return "o"
	case Cross:
		return "x"
	default:
		return "$"
	}
}

func (err InvalidMoveError) Error() string {
	return string(err)
}

func (board *Board) move(row int, column int, symbol Symbol) error {
	if row > 2 || column > 2 {
		return InvalidMoveError(fmt.Sprintf("Invalid index, %d:%d is not within the boundaries of the board", row, column))
	}
	if board[row][column] != None {
		return InvalidMoveError(fmt.Sprintf("Invalid move, %d:%d is being used by %s", row, column, (&symbol).Name()))
	}
	board[row][column] = symbol
	return nil
}

func (board *Board) toString() string {
	var res string
	rowSep := "-------\n"
	for i, row := range board {
		res += " "
		sep := "|"
		for j, symbol := range row {
			if j == 2 {
				sep = ""
			}
			res += symbol.String() + sep
		}
		if i == 2 {
			rowSep = ""
		}
		res += "\n" + rowSep
	}
	return res
}
