package models

import "time"

type Payment struct {
	ID          uint
	UserID      uint
	OrderID     string `gorm:"uniqueIndex"`
	PackageID   uint
	Amount      int64
	Status      string // pending, success, failed
	Type        string // pro, event
	ReferenceID uint   // event_id / subscription_id
	CreatedAt   time.Time
}
