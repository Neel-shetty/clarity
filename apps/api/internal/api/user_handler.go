package api

import (
	"net/http"

	"github.com/Neel-shetty/clarity/internal/domain"
	"github.com/Neel-shetty/clarity/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var user domain.SignUp
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Payload", "error message": err.Error(), "status": 0})
		return
	}
	err := h.userService.Signup(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to create the profile", "error": err.Error(), "status": 0})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created Successfully !!!", "status": 1})
}

func (h *UserHandler) Login(c *gin.Context) {
	var login domain.Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Payload", "error message": err.Error(), "status": 0})
		return
	}
	err = h.userService.Login(c.Request.Context(), &login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "error message": err.Error(), "status": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login Successfully", "status": 1})
}
