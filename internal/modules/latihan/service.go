package latihan

import (
	"errors"
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"time"
)

// ===== USER =====

// list latihan milik materi
func GetLatihanByMateri(materiID uint) ([]models.Latihan, error) {
	var list []models.Latihan
	err := database.DB.
		Where("materi_id = ?", materiID).
		Order("`order` ASC").
		Find(&list).Error
	return list, err
}

// ambil latihan by id
func GetLatihanByID(id uint) (*models.Latihan, error) {
	var latihan models.Latihan
	err := database.DB.First(&latihan, id).Error
	return &latihan, err
}

// ambil soal + opsi
func GetQuestionsByLatihan(latihanID uint) ([]models.Question, error) {
	var questions []models.Question
	err := database.DB.
		Where("latihan_id = ?", latihanID).
		Find(&questions).Error
	return questions, err
}

func GetOptionsByQuestion(questionID uint) ([]models.Option, error) {
	var options []models.Option
	err := database.DB.
		Where("question_id = ?", questionID).
		Find(&options).Error
	return options, err
}

// ===== ADMIN =====

func CreateLatihan(l *models.Latihan) error {
	return database.DB.Create(l).Error
}

func CreateQuestion(q *models.Question) error {
	return database.DB.Create(q).Error
}

func CreateOption(o *models.Option) error {
	return database.DB.Create(o).Error
}

func SubmitUserAnswer(
	userID uint,
	latihanID uint,
	attemptID uint,
	questionID uint,
	optionID uint,
) (bool, error) {

	// 1. CEK SUDAH JAWAB BELUM
	var existing models.UserAnswer
	if err := database.DB.
		Where(
			"user_id = ? AND question_id = ? AND attempt_id = ?",
			userID, questionID, attemptID,
		).
		First(&existing).Error; err == nil {
		return false, errors.New("already_answered")
	}

	// 2. VALIDASI OPTION MILIK QUESTION
	var option models.Option
	if err := database.DB.
		Where("id = ? AND question_id = ?", optionID, questionID).
		First(&option).Error; err != nil {
		return false, errors.New("invalid_option")
	}

	// 3. SIMPAN JAWABAN
	answer := models.UserAnswer{
		UserID:     userID,
		LatihanID:  latihanID,
		AttemptID:  attemptID,
		QuestionID: questionID,
		OptionID:   optionID,
		IsCorrect:  option.IsCorrect,
		AnsweredAt: time.Now(),
	}

	if err := database.DB.Create(&answer).Error; err != nil {
		return false, err
	}

	return option.IsCorrect, nil
}
