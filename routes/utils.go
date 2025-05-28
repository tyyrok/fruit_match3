package routes

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"sort"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"
)


func updateState(state *GameBoard, combs *[]Combination, t *Turn) *GameBoard {
	new_score := getScoresForCombs(combs)
	state.Scores += new_score
	newElems := getNewElems(new_score)
	if t != nil && *t != (Turn{}) {
		state.updateBoardByTurn(t)
	}
	state.updateBoard(combs, *newElems)
	return state	
}

func getNewElems(number int) *[]int {
	var res []int
	for i := 0; i < number; i++ {
		res = append(res, rand.IntN(4))
	}
	return &res
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
	sort.Slice(combs, func(i, j int) bool {
		return combs[i].getLenght() > combs[j].getLenght()
	})
	for _, val := range combs {
		toAdd := val
		for j := 0; j < len(combs); j ++ {
			if val.checkIntersection(&combs[j]) {
				if val.getLenght() < combs[j].getLenght() {
					toAdd = combs[j]
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

func getHighScores(pool *pgxpool.Pool) ([]HighScore, error) {
	rows, err := pool.Query(
		context.Background(),
		`SELECT m3tops.id, m3tops.scores FROM m3tops ORDER BY m3tops.scores DESC LIMIT 5;`)
	if err != nil {
		log.Printf("Failed to retrieve data %s", err)
		return nil, err
	}
	defer rows.Close()
	var highScores []HighScore
	for rows.Next() {
		var s HighScore
		err := rows.Scan(&s.Id, &s.Scores)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		highScores = append(highScores, s)
	}
	log.Println(highScores)
	return highScores, nil
}

func saveGameResult(state *GameBoard, ctx *gin.Context) {
	if state.Scores == 0 {
		return
	}
	dbctx, ok := ctx.Get("db")
	if !ok {
		log.Print("Failed to connect to db")
		return
	}
	pool, ok := dbctx.(*pgxpool.Pool)
	if !ok {
		log.Print("Failed to connect to db")
		return
	}
	_, err := pool.Exec(
			context.Background(),
			`INSERT INTO m3tops(scores) VALUES ($1) RETURNING id;`,
			state.Scores)
	if err != nil {
		log.Printf("Error while saving result %s\n", err)
	}
}