package routes

import (
	"net/http"
	"log"

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
	log.Println("Connects ws")
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Read msg error: %s", err)
			break
		}
		log.Printf("received: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Printf("Write msg error: %s", err)
			break
		}
	}
}