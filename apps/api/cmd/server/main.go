package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Neel-shetty/clarity/internal/api"
	"github.com/Neel-shetty/clarity/internal/config"
	"github.com/Neel-shetty/clarity/internal/middleware"
	"github.com/Neel-shetty/clarity/internal/repository"
	"github.com/Neel-shetty/clarity/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Health check working", "status": 1})
}

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}

	redisClient, err := config.ConnectRedisDB()
	if err != nil {
		fmt.Println("Redis client is not initialized")
	}
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService, redisClient)
	authMiddleware := middleware.AuthMiddleware(redisClient)

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "PATCH", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))
	r.GET("/health", health)
	api := r.Group("/api/v1")
	api.POST("/signup", userHandler.Signup)
	api.POST("/login", userHandler.Login)
	authorized := api.Group("")
	authorized.Use(authMiddleware)
	{
		authorized.GET("/profile", userHandler.GetProfile)
		authorized.POST("/logout", userHandler.Logout)
	}

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
