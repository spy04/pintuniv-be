package models

import "time"

type Quote struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Text     string `gorm:"type:text;not null" json:"text"`
	Author   string `gorm:"type:varchar(100)" json:"author"`
	IsActive bool   `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
