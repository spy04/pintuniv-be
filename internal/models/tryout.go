package models

import "time"

type Tryout struct {
	ID        uint
	Title     string
	Duration  int // menit
	IsActive  bool
	CreatedAt time.Time
}
