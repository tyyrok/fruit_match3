package routes

import (
	"net/http"
	"log"
	"encoding/json"
	"math/rand/v2"

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

	res := getStartGameBoard()
	data, _ := json.Marshal(res)
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
		log.Printf("received: %s with type: %d", message, mt)
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
		res, _ := processMessage(&msg)
		data, _ = json.Marshal(res)
		err = c.WriteMessage(mt, data)
		if err != nil {
			log.Printf("Write msg error: %s", err)
			break
		}
	}
}


func processMessage(msg *Message) (*Message, error) {
	return &Message{}, nil
}

func getStartGameBoard() *Message {
	var board GameBoard
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			board.Cells[i][j] = rand.IntN(4)
		}
	}
	return &Message{
		Type: "update_board",
		Data: map[string]any{"board": board.Cells},
	}
}
