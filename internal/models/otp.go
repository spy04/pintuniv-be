package models

import "time"

type OTP struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	CodeHash  string `gorm:"not null"`
	ExpiresAt time.Time
	Used      bool `gorm:"default:false"`
	CreatedAt time.Time
}
