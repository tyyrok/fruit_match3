package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// Set the router as the default one shipped with Gin
	router := gin.Default()
	expectedHost := os.Getenv("ORIGIN")

	// Setup Security Headers
	router.Use(func(c *gin.Context) {
		if c.Request.Host != expectedHost {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src 'self' 'unsafe-inline' 'unsafe-eval' *; script-src-elem * 'unsafe-inline' *; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})

	return router
}


func Run(dbpool *pgxpool.Pool) {
	httpPort := os.Getenv("PORT")
	router := NewRouter()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("db", dbpool)
		ctx.Next()
	})
	router.LoadHTMLGlob("templates/*")

	v1 := router.Group("/api", CheckOrigin())

	router.GET("/", mainPageHandler)

	v1.GET("/start", gameHandler)

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}


func CheckOrigin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		schema := os.Getenv("SCHEMA")
		expectedHost := os.Getenv("ORIGIN")
		host := fmt.Sprintf("%s%s", schema, expectedHost)

		origin := ctx.GetHeader("Origin")
		referer := ctx.GetHeader("Referer")

		if origin != "" && origin != host {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid origin"})
			return
		}
		if referer != "" && !strings.HasPrefix(referer, host) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid referer"})
			return
		}

		ctx.Next()
	}
}