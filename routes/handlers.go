package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func mainPageHandler(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"index.html",
		gin.H{"Name": "Gin Framework"})
}