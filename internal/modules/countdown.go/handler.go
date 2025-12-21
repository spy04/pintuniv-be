package countdown

import (
	"net/http"
	"time"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"

	"github.com/gin-gonic/gin"
)

func GetCountdown(c *gin.Context) {
	var countdown models.Countdown

	err := database.DB.
		Where("is_active = ?", true).
		Order("id desc").
		First(&countdown).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	now := time.Now()
	diff := countdown.TargetAt.Sub(now)

	if diff <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"is_started": true,
			"name":       countdown.Name,
			"target_at":  countdown.TargetAt,
		})
		return
	}

	days := int(diff.Hours()) / 24
	hours := int(diff.Hours()) % 24
	minutes := int(diff.Minutes()) % 60
	seconds := int(diff.Seconds()) % 60

	c.JSON(http.StatusOK, gin.H{
		"name":       countdown.Name,
		"target_at":  countdown.TargetAt,
		"now":        now,
		"days":       days,
		"hours":      hours,
		"minutes":    minutes,
		"seconds":    seconds,
		"is_started": false,
	})
}

type UpdateCountdownRequest struct {
	Name     string `json:"name"`
	TargetAt string `json:"target_at" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

func AdminUpsertCountdown(c *gin.Context) {
	var req UpdateCountdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	targetAt, err := time.Parse(time.RFC3339, req.TargetAt)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid target_at"})
		return
	}

	// nonaktifkan semua dulu
	database.DB.Model(&models.Countdown{}).
		Update("is_active", false)

	countdown := models.Countdown{
		Name:     req.Name,
		TargetAt: targetAt,
		IsActive: true,
	}

	database.DB.Create(&countdown)

	c.JSON(200, gin.H{
		"data": countdown,
	})
}
