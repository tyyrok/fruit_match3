package routes

import (
	"errors"
	"strconv"
)

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