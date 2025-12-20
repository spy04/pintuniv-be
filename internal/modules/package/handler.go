package package_module

import (
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreatePackageRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	DurationDay int    `json:"duration_day"`
}

func AdminCreatePackage(c *gin.Context) {
	var req CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pkg := models.Package{
		Code:        req.Code,
		Name:        req.Name,
		Price:       req.Price,
		DurationDay: req.DurationDay,
		IsActive:    true,
	}

	if err := database.DB.Create(&pkg).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_create_package"})
		return
	}

	c.JSON(201, pkg)
}

func AdminUpdatePackage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var pkg models.Package
	if err := database.DB.First(&pkg, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "package_not_found"})
		return
	}

	var req CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	database.DB.
		Model(&pkg).
		Updates(models.Package{
			Code:        req.Code,
			Name:        req.Name,
			Price:       req.Price,
			DurationDay: req.DurationDay,
		})

	c.JSON(200, pkg)
}

func AdminDeletePackage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	database.DB.
		Model(&models.Package{}).
		Where("id = ?", id).
		Update("is_active", false)

	c.JSON(200, gin.H{"status": "deleted"})
}

func ListPackages(c *gin.Context) {
	var packages []models.Package
	database.DB.
		Where("is_active = true").
		Find(&packages)

	c.JSON(200, packages)
}
