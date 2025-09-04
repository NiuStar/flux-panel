package handlers

import (
	"net/http"

	"github.com/flux-panel/go-backend/models"
	"github.com/flux-panel/go-backend/services"
	"github.com/gin-gonic/gin"
)

// NodeHandler wraps node related HTTP handlers.
type NodeHandler struct {
	Service *services.NodeService
}

// NewNodeHandler creates a new NodeHandler.
func NewNodeHandler(s *services.NodeService) *NodeHandler {
	return &NodeHandler{Service: s}
}

// RegisterRoutes registers node routes on the engine.
func (h *NodeHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/nodes", h.createNode)
	r.GET("/nodes", h.listNodes)
}

func (h *NodeHandler) createNode(c *gin.Context) {
	var node models.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateNode(&node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) listNodes(c *gin.Context) {
	nodes, err := h.Service.ListNodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}
