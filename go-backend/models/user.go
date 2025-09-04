package models

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User represents a system user.
type User struct {
	ID        uint   `gorm:"primaryKey"`
	User      string `gorm:"uniqueIndex;size:64"`
	Password  string `gorm:"size:128"`
	RoleID    int    `gorm:"column:role_id"`
	Status    int    `gorm:"column:status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// InitDB initializes a gorm.DB connection using the provided DSN.
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}, &Node{}, &Tunnel{}); err != nil {
		return nil, err
	}
	return db, nil
}
