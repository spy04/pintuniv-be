package models

import "time"

type Question struct {
	ID        uint
	LatihanID uint   // 👈 NEMPEL KE MATERI
	Question  string `gorm:"type:longtext"` // markdown
	CreatedAt time.Time
}
