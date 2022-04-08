package models

import (
	"time"

	"gorm.io/gorm"
)

type Node struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Address   string `json:"address" gorm:"index"`
	CreatedAt time.Time
}
