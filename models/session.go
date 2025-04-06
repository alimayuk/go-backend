package models

import "time"

type Session struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	TokenID   string `gorm:"unique"`
	UserAgent string
	IP        string
	Revoked   bool
	CreatedAt time.Time
}
