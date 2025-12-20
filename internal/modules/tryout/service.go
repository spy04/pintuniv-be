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

func GetRemainingSeconds(
	startedAt time.Time,
	durationMinutes int,
) int64 {
	endTime := startedAt.Add(
		time.Duration(durationMinutes) * time.Minute,
	)

	remaining := time.Until(endTime).Seconds()
	if remaining < 0 {
		return 0
	}
	return int64(remaining)
}
