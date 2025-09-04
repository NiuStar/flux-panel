package models

import "time"

// Node represents a server node that can host tunnels.
type Node struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64"`
	Address   string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
