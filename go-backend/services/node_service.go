package services

import (
	"github.com/flux-panel/go-backend/models"
	"gorm.io/gorm"
)

// NodeService contains business logic for nodes.
type NodeService struct {
	DB *gorm.DB
}

// NewNodeService creates a new NodeService.
func NewNodeService(db *gorm.DB) *NodeService {
	return &NodeService{DB: db}
}

// CreateNode inserts a new node record.
func (s *NodeService) CreateNode(n *models.Node) error {
	return s.DB.Create(n).Error
}

// ListNodes returns all nodes.
func (s *NodeService) ListNodes() ([]models.Node, error) {
	var nodes []models.Node
	if err := s.DB.Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}
