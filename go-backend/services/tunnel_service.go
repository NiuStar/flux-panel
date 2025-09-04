package services

import (
	"github.com/flux-panel/go-backend/models"
	"gorm.io/gorm"
)

// TunnelService contains business logic for tunnels.
type TunnelService struct {
	DB *gorm.DB
}

// NewTunnelService creates a new TunnelService.
func NewTunnelService(db *gorm.DB) *TunnelService {
	return &TunnelService{DB: db}
}

// CreateTunnel inserts a new tunnel record.
func (s *TunnelService) CreateTunnel(t *models.Tunnel) error {
	return s.DB.Create(t).Error
}

// ListTunnels returns all tunnels.
func (s *TunnelService) ListTunnels() ([]models.Tunnel, error) {
	var tunnels []models.Tunnel
	if err := s.DB.Find(&tunnels).Error; err != nil {
		return nil, err
	}
	return tunnels, nil
}
