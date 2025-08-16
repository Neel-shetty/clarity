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

	"github.com/Neel-shetty/clarity/internal/handler"
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

	appConfig := config.Load()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	redisSessionService := service.NewRedisSessionService(redisClient)
	userHandler := handler.NewUserHandler(userService, redisSessionService, appConfig)
	authMiddleware := middleware.AuthMiddleware(redisClient)

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	r.Use(cors.New(corsConfig))
	r.GET("/health", health)
	r.POST("/api/v1/signup", userHandler.Signup)
	r.POST("/api/v1/login", userHandler.Login)
	authorized := r.Group("/api/v1")
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
