package models

import (
	"time"

	"gorm.io/gorm"
)

type Business struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Category    string         `json:"category" gorm:"not null"`
	Latitude    float64        `json:"latitude" gorm:"not null"`
	Longitude   float64        `json:"longitude" gorm:"not null"`
	Address     string         `json:"address"`
	Phone       string         `json:"phone"`
	Website     string         `json:"website"`
	Description string         `json:"description"`
	Rating      float32        `json:"rating" gorm:"default:0"`
	ReviewCount int            `json:"review_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Reviews []Review `json:"reviews,omitempty" gorm:"foreignKey:BusinessID"`
}
