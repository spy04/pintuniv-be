package promo

import (
	"fmt"
	"os"
	"path"
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CreatePromoRequest struct {
	Code  string `json:"code" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"` // percent / flat
	Value int64  `json:"value" binding:"required"`

	StartAt string `json:"start_at" binding:"required"`
	EndAt   string `json:"end_at" binding:"required"`
}

func AdminCreatePromo(c *gin.Context) {
	var req CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	start, _ := time.Parse(time.RFC3339, req.StartAt)
	end, _ := time.Parse(time.RFC3339, req.EndAt)

	promo := models.Promo{
		Code:     req.Code,
		Name:     req.Name,
		Type:     req.Type,
		Value:    req.Value,
		StartAt:  start,
		EndAt:    end,
		IsActive: true,
	}

	database.DB.Create(&promo)

	c.JSON(201, promo)
}

type AttachPromoPackageRequest struct {
	PackageID uint `json:"package_id" binding:"required"`
}

func AdminAttachPromoPackage(c *gin.Context) {
	promoID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req AttachPromoPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	link := models.PromoPackage{
		PromoID:   uint(promoID),
		PackageID: req.PackageID,
	}

	database.DB.Create(&link)

	c.JSON(201, link)
}

func AdminUpdatePromo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var promo models.Promo
	if err := database.DB.First(&promo, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "promo_not_found"})
		return
	}

	var req CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	start, _ := time.Parse(time.RFC3339, req.StartAt)
	end, _ := time.Parse(time.RFC3339, req.EndAt)

	database.DB.Model(&promo).Updates(models.Promo{
		Code:    req.Code,
		Name:    req.Name,
		Type:    req.Type,
		Value:   req.Value,
		StartAt: start,
		EndAt:   end,
	})

	c.JSON(200, promo)
}
func AdminDisablePromo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	database.DB.
		Model(&models.Promo{}).
		Where("id = ?", id).
		Update("is_active", false)

	c.JSON(200, gin.H{"status": "disabled"})
}
func ValidatePromo(c *gin.Context) {
	code := c.Query("code")
	packageID, _ := strconv.ParseUint(c.Query("package_id"), 10, 64)

	if code == "" || packageID == 0 {
		c.JSON(400, gin.H{"valid": false, "reason": "invalid_param"})
		return
	}

	var promo models.Promo
	err := database.DB.
		Where(
			"code = ? AND is_active = true AND start_at <= NOW() AND end_at >= NOW()",
			code,
		).
		First(&promo).Error

	if err != nil {
		c.JSON(200, gin.H{"valid": false, "reason": "promo_not_found"})
		return
	}

	var count int64
	database.DB.
		Model(&models.PromoPackage{}).
		Where("promo_id = ? AND package_id = ?", promo.ID, packageID).
		Count(&count)

	if count == 0 {
		c.JSON(200, gin.H{"valid": false, "reason": "promo_not_applicable"})
		return
	}

	c.JSON(200, gin.H{
		"valid": true,
		"type":  promo.Type,
		"value": promo.Value,
	})
}

func AdminUploadPromoImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid_promo_id"})
		return
	}

	var promo models.Promo
	if err := database.DB.First(&promo, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "promo_not_found"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image_required"})
		return
	}

	if file.Size > 2*1024*1024 {
		c.JSON(400, gin.H{"error": "image_too_large"})
		return
	}

	ext := strings.ToLower(path.Ext(file.Filename))
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}
	if !allowed[ext] {
		c.JSON(400, gin.H{"error": "invalid_image_type"})
		return
	}

	uploadDir := "public/promos"
	_ = os.MkdirAll(uploadDir, os.ModePerm)

	filename := fmt.Sprintf("promo_%d_%d%s", promo.ID, time.Now().Unix(), ext)
	fullPath := path.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(500, gin.H{"error": "upload_failed"})
		return
	}

	// hapus file lama (optional tapi bagus)
	if promo.ImageURL != "" {
		oldPath := strings.TrimPrefix(promo.ImageURL, "/")
		_ = os.Remove(oldPath)
	}

	imageURL := fmt.Sprintf("/public/promos/%s", filename)
	database.DB.Model(&promo).Update("image_url", imageURL)

	c.JSON(200, gin.H{"image_url": imageURL})
}

func GetPromos(c *gin.Context) {
	var promos []models.Promo

	database.DB.
		Where(
			"is_active = true AND start_at <= NOW() AND end_at >= NOW()",
		).
		Order("created_at desc").
		Find(&promos)

	c.JSON(200, gin.H{
		"data": promos,
	})
}
