package quote

import (
	"math/rand"
	"net/http"
	"time"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"

	"github.com/gin-gonic/gin"
)

func GetQuote(c *gin.Context) {
	var quotes []models.Quote

	database.DB.
		Where("is_active = ?", true).
		Find(&quotes)

	if len(quotes) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	rand.Seed(time.Now().Unix())
	quote := quotes[rand.Intn(len(quotes))]

	c.JSON(http.StatusOK, gin.H{
		"data": quote,
	})
}

type CreateQuoteRequest struct {
	Text     string `json:"text" binding:"required"`
	Author   string `json:"author"`
	IsActive *bool  `json:"is_active"`
}

func AdminCreateQuote(c *gin.Context) {
	var req CreateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	quote := models.Quote{
		Text:     req.Text,
		Author:   req.Author,
		IsActive: isActive,
	}

	if err := database.DB.Create(&quote).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"data": quote})
}

func AdminUpdateQuote(c *gin.Context) {
	id := c.Param("id")

	var quote models.Quote
	if err := database.DB.First(&quote, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "quote not found"})
		return
	}

	var req CreateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	quote.Text = req.Text
	quote.Author = req.Author
	if req.IsActive != nil {
		quote.IsActive = *req.IsActive
	}

	database.DB.Save(&quote)
	c.JSON(200, gin.H{"data": quote})
}

func AdminDeleteQuote(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Quote{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}
