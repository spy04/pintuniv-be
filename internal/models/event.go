package models

import "time"

type Event struct {
	ID        uint
	Title     string
	StartAt   time.Time
	EndAt     time.Time
	IsActive  bool
	CreatedAt time.Time
}
