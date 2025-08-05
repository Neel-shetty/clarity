package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Unauthorized"})
			return
		}
		// sessionID := uuid.New().String()
		sessionKey := "session:" + sessionID
		userId, err := redisClient.Get(context.Background(), "session:"+sessionID).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Unauthorized"})
			return
		}
		err = redisClient.Expire(context.Background(), sessionKey, 24*time.Hour).Err()
		if err != nil {
			log.Printf("Could not refresh session for user %s: %v", userId, err)
		}
		c.Set("userId", userId)
		c.Next()

	}
}
