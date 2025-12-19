package models

type TryoutOption struct {
	ID         uint
	QuestionID uint
	Text       string `gorm:"type:text"`
	IsCorrect  bool
}
