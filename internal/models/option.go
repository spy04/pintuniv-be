package models

type Option struct {
	ID         uint
	QuestionID uint
	Text       string `gorm:"type:text"`
	IsCorrect  bool   `gorm:"default:false"`
}
