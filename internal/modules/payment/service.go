package payment

import (
	"time"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
)

func ExpireOldPendingPayments() error {
	expireBefore := time.Now().Add(-30 * time.Minute)

	return database.DB.
		Model(&models.Payment{}).
		Where(
			"status = ? AND created_at < ?",
			"pending",
			expireBefore,
		).
		Update("status", "failed").
		Error
}
