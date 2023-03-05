package logic

type Symbol byte

const (
	None Symbol = iota
	Circle
	Cross
)

type Board [3][3]Symbol

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

func (board *Board) ToString() string {
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

func (board *Board) ToByteArray() (array [9]byte) {
	for i, column := range board {
		for j, element := range column {
			array[i*3+j] = byte(element)
		}
	}
	return array
}

func BoardFromArray(array *[9]byte) (board Board) {
	for i := 0; i < 9; i++ {
		board[i/3][i%3] = Symbol((*array)[i])
	}
	return
}
