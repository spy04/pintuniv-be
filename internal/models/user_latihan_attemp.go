package models

import "time"

type UserLatihanAttempt struct {
	ID        uint
	UserID    uint
	LatihanID uint
	Score     int
	CreatedAt time.Time
}
