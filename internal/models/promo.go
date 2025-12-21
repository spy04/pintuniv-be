package models

import "time"

type Promo struct {
	ID uint

	Code string // PROMO20
	Name string

	Type  string // percent | flat
	Value int64  // 20 (percent) | 50000 (flat)

	StartAt time.Time
	EndAt   time.Time

	IsActive bool

	CreatedAt time.Time
}
