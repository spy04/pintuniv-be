package tryout

import (
	"time"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
)

func CountUserTryoutThisMonth(userID uint) (int64, error) {
	now := time.Now().UTC()

	startOfMonth := time.Date(
		now.Year(),
		now.Month(),
		1,
		0, 0, 0, 0,
		time.UTC,
	)

	var count int64
	err := database.DB.
		Model(&models.TryoutAttempt{}).
		Where(
			"user_id = ? AND started_at >= ?",
			userID, startOfMonth,
		).
		Count(&count).Error

	return count, err
}
