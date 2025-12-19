package models

import "time"

type TryoutAttempt struct {
	ID         uint
	UserID     uint
	TryoutID   uint
	StartedAt  time.Time
	FinishedAt *time.Time
	Score      int
}
