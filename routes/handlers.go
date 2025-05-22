package routes

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func mainPageHandler(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"index.html",
		gin.H{"Name": "Gin Framework"})
}

func gameHandler(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
		return
	}
	defer c.Close()

	gameState := getInitialGameState()
	data, _ := json.Marshal(&Message{
		Type: "update_board",
		Data: map[string]any{"board": gameState.Cells},
	})
	err = c.WriteMessage(1, data)
	if err != nil {
		log.Printf("Write msg error: %s", err)
		return
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Read msg error: %s", err)
			break
		}
		log.Printf("received: %s", message)
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			res, _ := json.Marshal("Wrong message type")
			err = c.WriteMessage(mt, res)
			if err != nil {
				log.Printf("Write msg error: %s", err)
				break
			}
		}
		res, _ := processMessage(&msg, gameState)
		ok := sendMessage(res, mt, c)
		if !ok {
			break
		}
		if res.Type == "end_game"{
			break
		}
		for {
			res, err := processAutoTurn(gameState)
			if err != nil {
				ok := sendMessage(&Message{Type: "resume"}, mt, c)
				if !ok {
					return
				}
				break
			}
			ok = sendMessage(res, mt, c)
			if !ok {
				return
			}
		}
	}
}

// Todo autoturn and board change
func processAutoTurn(state *GameBoard) (*Message, error) {
	combs := findCombinations(&state.Cells)
	if len(combs) > 0 {
		_ = updateState(state, &combs, &Turn{})
		return &Message{
			Type: "automove",
			Data: map[string]any{
				"status": "success",
				"turns": combs,
				"board": state.Cells,
				"scores": state.Scores},
		}, nil
	}
	return &Message{}, errors.New("not found automove")
}

func sendMessage(msg *Message, mt int, c *websocket.Conn) bool {
	data, _ := json.Marshal(msg)
	err := c.WriteMessage(mt, data)
	if err != nil {
		log.Printf("Write msg error: %s", err)
		return false
	}
	return true
}

func processMessage(msg *Message, state *GameBoard) (*Message, error) {
	switch msg.Type {
	case "move":
		return processTurn(msg, state)
	case "end_game":
		return processEndGame(state)
	default:
		return &Message{
			Type: "error",
		}, nil
	}
}

func processTurn(msg *Message, state *GameBoard) (*Message, error) {
	turn, err := validateTurn(msg)
	if err != nil {
		return &Message{Type: "turn validation error",}, nil
	}

	new_board := copyArray(state.Cells)
	new_board[turn.FromRow][turn.FromCol], new_board[turn.ToRow][turn.ToCol] = new_board[turn.ToRow][turn.ToCol], new_board[turn.FromRow][turn.FromCol]
	combs := findCombinations(new_board)
	if len(combs) > 0 {
		_ = updateState(state, &combs, turn)
		return &Message{
			Type: "move",
			Data: map[string]any{
				"status": "success",
				"turns": combs,
				"board": state.Cells,
				"scores": state.Scores,
				"turn": turn},
		}, nil
	}
	return &Message{
		Type: "move",
		Data: map[string]any{"status": "failure"},
	}, nil
}

func processEndGame(state *GameBoard) (*Message, error) {
	// TODO: implement saving game res
	return &Message{Type: "end_game", Data: map[string]any{"score": state.Scores}}, nil
}

func getInitialGameState() *GameBoard {
	var board GameBoard
	for {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				board.Cells[i][j] = rand.IntN(4)
			}
		}
		combs := findCombinations(&board.Cells)
		if len(combs) == 0 {
			break
		}
	}
	return &board
}
