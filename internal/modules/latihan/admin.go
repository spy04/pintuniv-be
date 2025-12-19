package latihan

import (
	// "net/http"
	"strconv"

	"pintuniv-go/internal/models"

	"github.com/gin-gonic/gin"
)

// =========================
// CREATE LATIHAN
// =========================
type CreateLatihanRequest struct {
	Title  string `json:"title" binding:"required"`
	IsFree bool   `json:"is_free"`
	Order  int    `json:"order"`
}

func AdminCreateLatihan(c *gin.Context) {
	materiID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CreateLatihanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	latihan := models.Latihan{
		MateriID: uint(materiID),
		Title:    req.Title,
		IsFree:   req.IsFree,
		Order:    req.Order,
	}

	if err := CreateLatihan(&latihan); err != nil {
		c.JSON(500, gin.H{"error": "failed to create latihan"})
		return
	}

	c.JSON(201, latihan)
}

// =========================
// CREATE QUESTION
// =========================
type CreateQuestionRequest struct {
	Question string `json:"question" binding:"required"`
}

func AdminCreateQuestion(c *gin.Context) {
	latihanID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	q := models.Question{
		LatihanID: uint(latihanID),
		Question:  req.Question,
	}

	CreateQuestion(&q)
	c.JSON(201, q)
}

// =========================
// CREATE OPTION
// =========================
type CreateOptionRequest struct {
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}

func AdminCreateOption(c *gin.Context) {
	questionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CreateOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	o := models.Option{
		QuestionID: uint(questionID),
		Text:       req.Text,
		IsCorrect:  req.IsCorrect,
	}

	CreateOption(&o)
	c.JSON(201, o)
}
