package promo

import (
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"strconv"
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
