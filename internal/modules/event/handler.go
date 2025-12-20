package event

import (
	"strconv"
	"time"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateEventRequest struct {
	Title   string `json:"title" binding:"required"`
	StartAt string `json:"start_at" binding:"required"` // ISO
	EndAt   string `json:"end_at" binding:"required"`
}

func AdminCreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	start, _ := time.Parse(time.RFC3339, req.StartAt)
	end, _ := time.Parse(time.RFC3339, req.EndAt)

	event := models.Event{
		Title:    req.Title,
		StartAt:  start,
		EndAt:    end,
		IsActive: true,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_create_event"})
		return
	}

	c.JSON(201, event)
}

type AttachTryoutRequest struct {
	TryoutID uint `json:"tryout_id" binding:"required"`
}

func AdminAttachTryout(c *gin.Context) {
	eventID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req AttachTryoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	link := models.EventTryout{
		EventID:  uint(eventID),
		TryoutID: req.TryoutID,
	}

	if err := database.DB.Create(&link).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_attach_tryout"})
		return
	}

	c.JSON(201, link)
}

func ListActiveEvents(c *gin.Context) {
	now := time.Now()

	var events []models.Event
	database.DB.
		Where(
			"is_active = true AND start_at <= ? AND end_at >= ?",
			now, now,
		).
		Find(&events)

	c.JSON(200, events)
}

func ListEventTryouts(c *gin.Context) {
	eventID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var tryouts []models.Tryout
	database.DB.
		Table("tryouts").
		Joins("JOIN event_tryouts et ON et.tryout_id = tryouts.id").
		Where("et.event_id = ?", eventID).
		Find(&tryouts)

	c.JSON(200, tryouts)
}
