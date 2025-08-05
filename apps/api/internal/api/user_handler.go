package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Neel-shetty/clarity/internal/domain"
	"github.com/Neel-shetty/clarity/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	userService service.UserService
	redisClient *redis.Client
}

func NewUserHandler(userService service.UserService, redisClient *redis.Client) *UserHandler {
	return &UserHandler{
		userService: userService,
		redisClient: redisClient,
	}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var user domain.SignUp
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid Request Payload",
				"error message": err.Error(),
				"status":        0})
		return
	}
	err := h.userService.Signup(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "Unable to create the profile !!!",
				"error":  err.Error(),
				"status": 0})
		return
	}
	c.JSON(http.StatusCreated,
		gin.H{"message": "User created Successfully !!!",
			"status": 1})
}

func (h *UserHandler) Login(c *gin.Context) {

	var login domain.Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid Request Payload",
				"error message": err.Error(),
				"status":        0})
		return
	}
	user, err := h.userService.Login(c.Request.Context(), &login)
	if err != nil {
		c.JSON(http.StatusUnauthorized,
			gin.H{"error": "Invalid email or password",
				"error message": err.Error(),
				"status":        0})
		return
	}
	sessionID := uuid.New().String()

	err = h.redisClient.Set(c.Request.Context(), "session:"+sessionID, user.Id.String(), 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to create the session"})
		return
	}

	c.SetCookie(
		"session_id",
		sessionID,
		3600*24,
		"/",
		"localhost",
		false,
		true,
	)
	c.JSON(http.StatusOK,
		gin.H{"message": "Login Successfully",
			"status": 1,
			"user":   user})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized,
			gin.H{"error": "User Id is not found in the session",
				"status": 0})
		return
	}
	fmt.Println(userId)
	userID, err := uuid.Parse(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format", "user": userID})
		return
	}

	profile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	user := domain.Profile{
		Name:  profile.Name,
		Email: profile.Email,
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusOK,
			gin.H{"message": "Already logged out",
				"status": 1})
		return
	}
	h.redisClient.Del(context.Background(), "session:"+sessionID)
	c.SetCookie("session_id", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK,
		gin.H{"message": "Logout Successfully",
			"status": 1})

}
