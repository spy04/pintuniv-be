package models

type TryoutQuestion struct {
	ID       uint
	TryoutID uint
	Question string `gorm:"type:longtext"`
	Order    int
}
