package api

import (
	"net/http"

	"github.com/HarshithRajesh/clarity/internal/domain"
	"github.com/HarshithRajesh/clarity/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var user domain.SignUp
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Payload", "error message": err.Error()})
		return
	}
	err := h.userService.SignUp(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to create the profile", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created Successfully !!!"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var login domain.Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Payload", "error message": err.Error()})
		return
	}
	err = h.userService.Login(c.Request.Context(), &login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "error message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login Successfully"})
}
