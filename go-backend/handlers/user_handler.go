package handlers

import (
	"net/http"

	"github.com/flux-panel/go-backend/models"
	"github.com/flux-panel/go-backend/services"
	"github.com/gin-gonic/gin"
)

// UserHandler wraps user related HTTP handlers.
type UserHandler struct {
	Service *services.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// RegisterRoutes registers routes on the given engine.
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/login", h.login)
	r.POST("/users", h.createUser)
	r.GET("/users", h.listUsers)
}

func (h *UserHandler) login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) listUsers(c *gin.Context) {
	users, err := h.Service.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
