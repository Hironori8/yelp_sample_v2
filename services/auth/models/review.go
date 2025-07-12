package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	BusinessID uint      `json:"business_id" gorm:"not null"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Rating     int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	
	Business Business `json:"business,omitempty" gorm:"foreignKey:BusinessID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}