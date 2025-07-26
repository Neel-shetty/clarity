package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HarshithRajesh/clarity/internal/config"
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Health check working"})
}

func main() {
	_, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}

	r := gin.Default()
	r.GET("/health", health)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	service := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen :%s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := service.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}
	log.Println("Server Exiting")
}
