package handlers

import (
	"net/http"

	"github.com/flux-panel/go-backend/models"
	"github.com/flux-panel/go-backend/services"
	"github.com/gin-gonic/gin"
)

// TunnelHandler wraps tunnel related HTTP handlers.
type TunnelHandler struct {
	Service *services.TunnelService
}

// NewTunnelHandler creates a new TunnelHandler.
func NewTunnelHandler(s *services.TunnelService) *TunnelHandler {
	return &TunnelHandler{Service: s}
}

// RegisterRoutes registers tunnel routes on the engine.
func (h *TunnelHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/tunnels", h.createTunnel)
	r.GET("/tunnels", h.listTunnels)
}

func (h *TunnelHandler) createTunnel(c *gin.Context) {
	var t models.Tunnel
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateTunnel(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TunnelHandler) listTunnels(c *gin.Context) {
	tunnels, err := h.Service.ListTunnels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tunnels)
}
