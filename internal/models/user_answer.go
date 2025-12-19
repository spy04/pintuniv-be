package models

import "time"

type UserAnswer struct {
	ID         uint
	UserID     uint
	LatihanID  uint
	QuestionID uint
	OptionID   uint
	AttemptID  uint
	IsCorrect  bool
	AnsweredAt time.Time
}
