package models

import "time"

// Tunnel represents a proxy tunnel associated with a node.
type Tunnel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64"`
	NodeID    uint   `gorm:"column:node_id"`
	Config    string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
