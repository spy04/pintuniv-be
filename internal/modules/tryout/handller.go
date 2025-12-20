package tryout

import (
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/modules/materi"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type CreateTryoutRequest struct {
	Title    string `json:"title" binding:"required"`
	Duration int    `json:"duration" binding:"required"` // menit
}

func AdminCreateTryout(c *gin.Context) {
	var req CreateTryoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tryout := models.Tryout{
		Title:    req.Title,
		Duration: req.Duration,
		IsActive: true,
	}

	if err := database.DB.Create(&tryout).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_create_tryout"})
		return
	}

	c.JSON(201, tryout)
}

type CreateTryoutQuestionRequest struct {
	Question string `json:"question" binding:"required"`
	Order    int    `json:"order"`
}

func AdminCreateTryoutQuestion(c *gin.Context) {
	tryoutID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CreateTryoutQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	q := models.TryoutQuestion{
		TryoutID: uint(tryoutID),
		Question: req.Question,
		Order:    req.Order,
	}

	if err := database.DB.Create(&q).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_create_question"})
		return
	}

	c.JSON(201, q)
}

type CreateTryoutOptionRequest struct {
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}

func AdminCreateTryoutOption(c *gin.Context) {
	questionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CreateTryoutOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	o := models.TryoutOption{
		QuestionID: uint(questionID),
		Text:       req.Text,
		IsCorrect:  req.IsCorrect,
	}

	if err := database.DB.Create(&o).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_create_option"})
		return
	}

	c.JSON(201, o)
}

func StartTryout(c *gin.Context) {
	userID := c.GetUint("user_id")
	tryoutID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// 1. CEK TRYOUT
	var t models.Tryout
	if err := database.DB.
		Where("id = ? AND is_active = true", tryoutID).
		First(&t).Error; err != nil {
		c.JSON(404, gin.H{"error": "tryout_not_found"})
		return
	}

	// 2. CEK USER PRO (PAKAI YANG SUDAH ADA)
	isPro := materi.IsUserPro(&userID)

	if !isPro {
		count, err := CountUserTryoutThisMonth(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed_check_limit"})
			return
		}

		if count >= 3 {
			c.JSON(403, gin.H{
				"error": "tryout_limit_reached",
				"limit": 3,
			})
			return
		}
	}

	// 3. BUAT ATTEMPT
	attempt := models.TryoutAttempt{
		UserID:    userID,
		TryoutID:  uint(tryoutID),
		StartedAt: time.Now(),
	}

	if err := database.DB.Create(&attempt).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_start_tryout"})
		return
	}

	remaining := GetRemainingSeconds(
		attempt.StartedAt,
		t.Duration,
	)

	c.JSON(201, gin.H{
		"attempt_id":        attempt.ID,
		"remaining_seconds": remaining,
	})

}

func GetTryoutQuestionByNumber(c *gin.Context) {
	// userID := c.GetUint("user_id")

	tryoutID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	number, _ := strconv.Atoi(c.Param("number"))

	attemptIDStr := c.Query("attempt_id")
	if attemptIDStr == "" {
		c.JSON(400, gin.H{"error": "attempt_id_required"})
		return
	}
	attemptID, _ := strconv.ParseUint(attemptIDStr, 10, 64)

	if number <= 0 {
		c.JSON(400, gin.H{"error": "invalid_question_number"})
		return
	}

	offset := number - 1

	// 1. AMBIL QUESTION KE-N
	var question models.TryoutQuestion
	err := database.DB.
		Where("tryout_id = ?", tryoutID).
		Order("`order` ASC, id ASC").
		Offset(offset).
		Limit(1).
		First(&question).Error

	if err != nil {
		c.JSON(404, gin.H{"error": "question_not_found"})
		return
	}

	// 2. AMBIL OPTIONS
	var options []models.TryoutOption
	database.DB.
		Where("question_id = ?", question.ID).
		Find(&options)

	// 3. CEK JAWABAN USER (UNTUK RESUME)
	var answer models.TryoutAnswer
	answered := true
	if err := database.DB.
		Where(
			"attempt_id = ? AND question_id = ?",
			attemptID, question.ID,
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

type SubmitTryoutRequest struct {
	AttemptID uint `json:"attempt_id" binding:"required"`
}

func SubmitTryout(c *gin.Context) {
	userID := c.GetUint("user_id")
	tryoutID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req SubmitTryoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// =========================
	// START TRANSACTION
	// =========================
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(500, gin.H{"error": "failed_start_transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// =========================
	// LOCK ATTEMPT ROW
	// =========================
	var attempt models.TryoutAttempt
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(
			"id = ? AND user_id = ? AND tryout_id = ?",
			req.AttemptID, userID, tryoutID,
		).
		First(&attempt).Error; err != nil {

		tx.Rollback()
		c.JSON(404, gin.H{"error": "attempt_not_found"})
		return
	}

	// =========================
	// SUDAH FINISH?
	// =========================
	if attempt.FinishedAt != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": "tryout_already_submitted"})
		return
	}

	// =========================
	// AMBIL TRYOUT (TANPA LOCK)
	// =========================
	var t models.Tryout
	if err := tx.First(&t, attempt.TryoutID).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "tryout_not_found"})
		return
	}

	// =========================
	// CEK WAKTU
	// =========================
	remaining := GetRemainingSeconds(
		attempt.StartedAt,
		t.Duration,
	)

	if remaining <= 0 {
		now := time.Now()
		tx.
			Model(&models.TryoutAttempt{}).
			Where("id = ?", attempt.ID).
			Update("finished_at", &now)

		tx.Commit()
		c.JSON(400, gin.H{"error": "time_up"})
		return
	}

	// =========================
	// AMBIL JAWABAN DRAFT
	// =========================
	var answers []models.TryoutAnswer
	tx.
		Where("attempt_id = ?", attempt.ID).
		Find(&answers)

	score := 0
	for _, a := range answers {
		var option models.TryoutOption
		if err := tx.First(&option, a.OptionID).Error; err == nil {
			if option.IsCorrect {
				score++
			}
		}
	}

	// =========================
	// FINALIZE ATTEMPT
	// =========================
	now := time.Now()
	if err := tx.
		Model(&models.TryoutAttempt{}).
		Where("id = ?", attempt.ID).
		Updates(map[string]interface{}{
			"score":       score,
			"finished_at": &now,
		}).Error; err != nil {

		tx.Rollback()
		c.JSON(500, gin.H{"error": "failed_finalize"})
		return
	}

	// =========================
	// COMMIT
	// =========================
	tx.Commit()

	c.JSON(200, gin.H{
		"score": score,
	})
}

func ResumeTryout(c *gin.Context) {
	userID := c.GetUint("user_id")
	tryoutID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var attempt models.TryoutAttempt
	err := database.DB.
		Where(
			"user_id = ? AND tryout_id = ? AND finished_at IS NULL",
			userID, tryoutID,
		).
		Order("started_at DESC").
		First(&attempt).Error

	if err != nil {
		// tidak ada attempt aktif
		c.JSON(200, gin.H{
			"attempt_id": nil,
		})
		return
	}

	// 👇 TAMBAHAN DI SINI (HITUNG SOAL TERAKHIR)
	var answeredCount int64
	database.DB.
		Model(&models.TryoutAnswer{}).
		Where("attempt_id = ?", attempt.ID).
		Count(&answeredCount)

	var t models.Tryout
	database.DB.First(&t, attempt.TryoutID)

	remaining := GetRemainingSeconds(
		attempt.StartedAt,
		t.Duration,
	)

	c.JSON(200, gin.H{
		"attempt_id":        attempt.ID,
		"next_number":       answeredCount + 1,
		"remaining_seconds": remaining,
	})

}

type SaveTryoutAnswerRequest struct {
	AttemptID  uint `json:"attempt_id" binding:"required"`
	QuestionID uint `json:"question_id" binding:"required"`
	OptionID   uint `json:"option_id" binding:"required"`
}

func SaveTryoutAnswer(c *gin.Context) {
	var req SaveTryoutAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// upsert: kalau sudah ada → update
	var answer models.TryoutAnswer
	err := database.DB.
		Where(
			"attempt_id = ? AND question_id = ?",
			req.AttemptID, req.QuestionID,
		).
		First(&answer).Error

	if err == nil {
		// update
		database.DB.
			Model(&answer).
			Update("option_id", req.OptionID)
	} else {
		// create
		answer = models.TryoutAnswer{
			AttemptID:  req.AttemptID,
			QuestionID: req.QuestionID,
			OptionID:   req.OptionID,
		}
		database.DB.Create(&answer)
	}

	c.JSON(200, gin.H{"status": "saved"})
}

func GetTryoutLeaderboard(c *gin.Context) {
	tryoutID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	type Row struct {
		UserID uint
		Name   string
		Score  int
	}

	// ambil attempt terakhir tiap user
	rows := []Row{}

	database.DB.Raw(`
		SELECT t1.user_id, t1.score
		FROM tryout_attempts t1
		INNER JOIN (
			SELECT user_id, MAX(started_at) AS last_attempt
			FROM tryout_attempts
			WHERE tryout_id = ? AND finished_at IS NOT NULL
			GROUP BY user_id
		) t2
		ON t1.user_id = t2.user_id AND t1.started_at = t2.last_attempt
		WHERE t1.tryout_id = ?
		ORDER BY t1.score DESC
		LIMIT ?
	`, tryoutID, tryoutID, limit).
		Scan(&rows)

	resp := []gin.H{}
	rank := 1

	for _, r := range rows {
		resp = append(resp, gin.H{
			"rank":    rank,
			"user_id": r.UserID,
			"score":   r.Score,
		})
		rank++
	}

	c.JSON(200, gin.H{
		"tryout_id":   tryoutID,
		"leaderboard": resp,
	})
}
