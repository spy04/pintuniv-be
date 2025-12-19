package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	Password     string `gorm:"not null"`
	RefreshToken string `gorm:"type:text"`

	IsActive   bool `gorm:"default:true"`
	IsVerified bool `gorm:"default:false"`

	Profile   Profile `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
