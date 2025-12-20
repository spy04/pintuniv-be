package models

import "time"

type Package struct {
	ID          uint
	Code        string // ex: PRO_MONTHLY
	Name        string // ex: Pro Bulanan
	Price       int64  // dalam rupiah
	DurationDay int    // 30, 365, dll (0 kalau one-time)
	IsActive    bool
	CreatedAt   time.Time
}
