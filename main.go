package main

import (
	"match3/routes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Println("Error loading .env file")
	}
	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)

	encodedPassword := os.Getenv("DB_PASSWORD")
	dbUrl := os.Getenv("DB_URL")
	connString := fmt.Sprintf(dbUrl, url.QueryEscape(encodedPassword))
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}
	defer dbpool.Close()

	InitDB(dbpool)
	routes.Run(dbpool)
}
