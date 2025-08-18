package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Neel-shetty/clarity/internal/config"
	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/Neel-shetty/clarity/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type UserHandler struct {
	userService    service.UserService
	sessionService service.SessionService
	cfg            config.Config
}

func NewUserHandler(userService service.UserService, sessionService service.SessionService, cfg config.Config) *UserHandler {
	return &UserHandler{
		userService:    userService,
		sessionService: sessionService,
		cfg:            cfg,
	}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var user model.SignUp
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid Request Payload",
				"error message": err.Error()})
		return
	}

	err := h.userService.Signup(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "Unable to create the profile !!!",
				"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated,
		gin.H{"message": "User created Successfully !!!"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var login model.Login

	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid Request Payload",
				"error message": err.Error()})
		return
	}

	user, err := h.userService.Login(c.Request.Context(), &login)
	if err != nil {
		c.JSON(http.StatusUnauthorized,
			gin.H{"error": "Invalid email or password",
				"error message": err.Error()})
		return
	}

	sessionID, err := h.sessionService.CreateSession(c.Request.Context(), user.ID, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to create the session"})
		return
	}

	c.SetCookie(
		"session_id",
		sessionID,
		int(h.cfg.CookieMaxAge),
		"/",
		h.cfg.CookieDomain,
		h.cfg.CookieSecure,
		true,
	)

	c.JSON(http.StatusOK,
		gin.H{"message": "Login Successfully"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized,
			gin.H{"error": "User Id is not found in the session"})
		return
	}

	userID, err := uuid.FromString(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	profile := model.Profile{
		FullName: user.Name,
	}
	log.Println(profile.ID.String())

	c.JSON(http.StatusOK, gin.H{
		"name": profile.FullName,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusOK,
			gin.H{"message": "Already logged out"})
		return
	}

	if err := h.sessionService.DeleteSession(context.Background(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Failed to logout: could not delete session",
				"error": err.Error()})
		return
	}

	c.SetCookie("session_id", "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

	c.JSON(http.StatusOK,
		gin.H{"message": "Logout Successfully"})

}
