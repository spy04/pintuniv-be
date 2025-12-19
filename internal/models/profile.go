package models

import "time"

type Profile struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"uniqueIndex"`
	FullName  string
	AvatarURL string `gorm:"type:varchar(255)"`
	Phone     string
	IsPro     bool `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
