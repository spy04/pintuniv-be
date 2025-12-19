package models

import "time"

type PasswordReset struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	TokenHash string `gorm:"not null"`
	ExpiresAt time.Time
	Used      bool `gorm:"default:false"`
	CreatedAt time.Time
}
