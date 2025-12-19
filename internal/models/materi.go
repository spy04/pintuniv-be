package models

import "time"

type Materi struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(255)"`
	Slug        string `gorm:"uniqueIndex"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
}
