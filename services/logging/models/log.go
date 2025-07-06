package models

import (
	"time"
)

type ReviewViewLog struct {
	UserID     int       `db:"user_id" json:"user_id"`
	BusinessID int       `db:"business_id" json:"business_id"`
	ReviewID   int       `db:"review_id" json:"review_id"`
	ViewedAt   time.Time `db:"viewed_at" json:"viewed_at"`
	IPAddress  string    `db:"ip_address" json:"ip_address"`
	UserAgent  string    `db:"user_agent" json:"user_agent"`
}

type LogRequest struct {
	UserID     int    `json:"user_id"` // PostgreSQLのusers.idと一致する整数型
	BusinessID int    `json:"business_id"`
	ReviewID   int    `json:"review_id"`
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
}

type LogResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
