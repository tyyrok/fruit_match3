package routes

import (
	"errors"
	"log"
	"strconv"
)


func updateState(state *GameBoard, combs *[]Combination, t *Turn) *GameBoard {
	state.Scores += getScoresForCombs(combs)
	//state.Cells = updateGameBoard(state.Cells, t)
	return state	
}

func updateGameBoard(b [8][8]int, t *Turn) [8][8]int {
	b[t.FromRow][t.FromCol], b[t.ToRow][t.ToCol] = b[t.ToRow][t.ToCol], b[t.FromRow][t.FromCol]
	return b
}

func getScoresForCombs(combs *[]Combination) int {
	scores := 0
	for _, comb := range *combs {
		scores += comb.getLenght()
	}
	return scores
}

func findCombinations(board *[8][8]int) []Combination {
	var combs []Combination
	combs = findVerticalCombs(board)
	combs = append(combs, findHorizontalCombs(board)...)
	combs = removeDuplicateCombs(combs)
	return combs
}

func removeDuplicateCombs(combs []Combination) []Combination {
	var res []Combination
	for _, val := range combs {
		toAdd := val
		for _, val_2 := range combs {
			if val.checkIntersection(&val_2) {
				if val.getLenght() < val_2.getLenght() {
					toAdd = val_2
				}
			}
		}
		if !checkIsAlreadyInSlice(res, toAdd) {
			res = append(res, toAdd)
		}
	}
	return res
}

func checkIsAlreadyInSlice(slice []Combination, elem Combination) bool {
	for _, e := range slice {
		if elem.equal(&e) {
			return true
		}
	}
	return false
}


func findHorizontalCombs(board *[8][8]int) []Combination {
	var combs []Combination
	for i := 0; i < 8; i ++ {
		for j := 0; j < 6; j ++ {
			if board[i][j] == board[i][j+1] && board[i][j] == board[i][j+2] {
				k := j + 1
				for k < 8 {
					if board[i][j] == board[i][k] {
						k += 1
					} else {
						break
					}
				}
				var points []Point
				for t := j; t < k; t ++ {
					points = append(points, Point{X: t, Y: i})
				}
				comb := Combination{Points: points}
				combs = append(combs, comb)
			}
		}
	}
	log.Println("Horizontal")
	log.Println(combs)
	return combs
}

func findVerticalCombs(board *[8][8]int) []Combination {
	var combs []Combination
	for i := 0; i < 8; i ++ {
		for j := 0; j < 6; j ++ {
			if board[j][i] == board[j+1][i] && board[j][i] == board[j+2][i] {
				k := j + 1
				for k < 8 {
					if board[j][i] == board[k][i] {
						k += 1
					} else {
						break
					}
				}
				var points []Point
				for t := j; t < k; t ++ {
					points = append(points, Point{X: i, Y: t})
				}
				comb := Combination{Points: points}
				combs = append(combs, comb)
			}
		}
	}
	log.Println("Vertical")
	log.Println(board)
	log.Println(combs)
	return combs
}

func validateTurn(msg *Message) (*Turn, error) {
	var fail bool
	from_row, ok := msg.Data["from_row"]
	if !ok {
		fail = true
	}
	f_row, err := toInt(from_row)
	if err != nil {
		fail = true
	}
	to_row, ok := msg.Data["to_row"]
	if !ok {
		fail = true
	}
	t_row, err := toInt(to_row)
	if err != nil {
		fail = true
	}
	from_col, ok := msg.Data["from_col"]
	if !ok {
		fail = true
	}
	f_col, err := toInt(from_col)
	if err != nil {
		fail = true
	}
	to_col, ok := msg.Data["to_col"]
	if !ok {
		fail = true
	}
	t_col, err := toInt(to_col)
	if err != nil {
		fail = true
	}
	if fail{
		return &Turn{}, errors.New("turn validation error")
	}
	return &Turn{FromRow: f_row, FromCol: f_col, ToRow: t_row, ToCol: t_col,}, nil
}

func toInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		r, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return r, nil
	default:
		return 0, errors.New("undefined type")
	}
}

func copyArray(arr [8][8]int) *[8][8]int {
	var dst [8][8]int
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			dst[i][j] = arr[i][j]
		}
	}
	return &dst
}