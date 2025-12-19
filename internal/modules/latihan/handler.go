package latihan

import (
	"net/http"
	"strconv"
	"time"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/modules/materi"

	"github.com/gin-gonic/gin"
)

// =========================
// GET /materi/:id/latihan
// =========================
func ListLatihanByMateri(c *gin.Context) {
	materiID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	list, err := GetLatihanByMateri(uint(materiID))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to load latihan"})
		return
	}

	var userID *uint = nil
	if uid, ok := c.Get("user_id"); ok {
		id := uid.(uint)
		userID = &id
	}
	isPro := materi.IsUserPro(userID)

	resp := []gin.H{}
	for _, l := range list {
		locked := !l.IsFree && !isPro

		resp = append(resp, gin.H{
			"id":      l.ID,
			"title":   l.Title,
			"is_free": l.IsFree,
			"locked":  locked,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// =========================
// GET /latihan/:id
// =========================
func GetLatihanDetail(c *gin.Context) {
	latihanID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	latihan, err := GetLatihanByID(uint(latihanID))
	if err != nil {
		c.JSON(404, gin.H{"error": "latihan not found"})
		return
	}

	var userID *uint = nil
	if uid, ok := c.Get("user_id"); ok {
		id := uid.(uint)
		userID = &id
	}
	isPro := materi.IsUserPro(userID)

	// 🔒 LOCK DI SINI
	if !latihan.IsFree && !isPro {
		c.JSON(403, gin.H{
			"error":  "latihan_locked",
			"reason": "pro_required",
		})
		return
	}

	questions, _ := GetQuestionsByLatihan(latihan.ID)

	respQuestions := []gin.H{}
	for _, q := range questions {
		options, _ := GetOptionsByQuestion(q.ID)

		opts := []gin.H{}
		for _, o := range options {
			opts = append(opts, gin.H{
				"id":   o.ID,
				"text": o.Text,
			})
		}

		respQuestions = append(respQuestions, gin.H{
			"id":       q.ID,
			"question": q.Question,
			"options":  opts,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"latihan": gin.H{
			"id":    latihan.ID,
			"title": latihan.Title,
		},
		"questions": respQuestions,
	})
}

type SubmitAnswerRequest struct {
	AttemptID uint `json:"attempt_id" binding:"required"`
	OptionID  uint `json:"option_id" binding:"required"`
}

func SubmitAnswer(c *gin.Context) {
	userID := c.GetUint("user_id")
	questionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// cari latihan_id dari question
	var question models.Question
	if err := database.DB.First(&question, questionID).Error; err != nil {
		c.JSON(404, gin.H{"error": "question not found"})
		return
	}

	isCorrect, err := SubmitUserAnswer(
		userID,
		question.LatihanID,
		req.AttemptID,
		uint(questionID),
		req.OptionID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to submit answer"})
		return
	}

	c.JSON(200, gin.H{
		"correct": isCorrect,
	})
}

func StartLatihan(c *gin.Context) {
	userID := c.GetUint("user_id")
	latihanID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	attempt := models.UserLatihanAttempt{
		UserID:    userID,
		LatihanID: uint(latihanID),
		Score:     0,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&attempt).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_start_attempt"})
		return
	}

	c.JSON(201, gin.H{
		"attempt_id": attempt.ID,
	})
}

func FinishLatihan(c *gin.Context) {
	userID := c.GetUint("user_id")
	attemptID, _ := strconv.ParseUint(c.Param("attempt_id"), 10, 64)

	var count int64
	database.DB.
		Model(&models.UserAnswer{}).
		Where(
			"user_id = ? AND attempt_id = ? AND is_correct = true",
			userID, attemptID,
		).
		Count(&count)

	database.DB.
		Model(&models.UserLatihanAttempt{}).
		Where("id = ?", attemptID).
		Update("score", int(count))

	c.JSON(200, gin.H{
		"score": count,
	})
}

func GetQuestionByNumber(c *gin.Context) {
	userID := c.GetUint("user_id")

	latihanID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	number, _ := strconv.Atoi(c.Param("number"))

	attemptIDStr := c.Query("attempt_id")
	if attemptIDStr == "" {
		c.JSON(400, gin.H{"error": "attempt_id_required"})
		return
	}
	attemptID, _ := strconv.ParseUint(attemptIDStr, 10, 64)

	// offset = nomor - 1
	offset := number - 1
	if offset < 0 {
		c.JSON(400, gin.H{"error": "invalid_number"})
		return
	}

	// ambil question ke-n
	var question models.Question
	err := database.DB.
		Where("latihan_id = ?", latihanID).
		Order("id ASC").
		Offset(offset).
		Limit(1).
		First(&question).Error

	if err != nil {
		c.JSON(404, gin.H{"error": "question_not_found"})
		return
	}

	// ambil options
	var options []models.Option
	database.DB.
		Where("question_id = ?", question.ID).
		Find(&options)

	// cek apakah user sudah jawab di attempt ini
	var answer models.UserAnswer
	answered := true
	if err := database.DB.
		Where(
			"user_id = ? AND question_id = ? AND attempt_id = ?",
			userID, question.ID, attemptID,
		).
		First(&answer).Error; err != nil {
		answered = false
	}

	respOptions := []gin.H{}
	for _, o := range options {
		respOptions = append(respOptions, gin.H{
			"id":   o.ID,
			"text": o.Text,
		})
	}

	c.JSON(200, gin.H{
		"question": gin.H{
			"id":       question.ID,
			"number":   number,
			"content":  question.Question,
			"answered": answered,
			"answer": gin.H{
				"option_id": answer.OptionID,
			},
		},
		"options": respOptions,
	})
}
