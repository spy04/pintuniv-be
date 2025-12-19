package tryout

import (
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
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
