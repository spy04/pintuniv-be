package models

import "time"

type Profile struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex;not null"`

	FullName    string `gorm:"type:varchar(100);not null"`
	AvatarURL   string `gorm:"type:varchar(255)"`
	Phone       string `gorm:"type:varchar(20)"`
	Bio         string `gorm:"type:varchar(500)"`
	AsalSekolah string `gorm:"type:varchar(100)"`
	Jurusan     string `gorm:"type:varchar(100)"`
	Lulus       int    `gorm:"check:lulus >= 1900 AND lulus <= 2100"`

	IsPro bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
