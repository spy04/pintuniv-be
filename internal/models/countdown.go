package models

import "time"

type Countdown struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(100)" json:"name"` // contoh: UTBK 2026
	TargetAt time.Time `json:"target_at"`
	IsActive bool      `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
