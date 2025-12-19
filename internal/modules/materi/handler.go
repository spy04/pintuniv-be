package materi

import (
	"fmt"
	"net/http"
	"path/filepath"
	"pintuniv-go/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ========================
// GET /api/materi
// ========================
func GetMateriList(c *gin.Context) {
	list, err := GetAllMateri()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch materi",
		})
		return
	}

	// materi TIDAK PERNAH di-lock
	// FE hanya butuh metadata
	resp := []gin.H{}

	for _, m := range list {
		resp = append(resp, gin.H{
			"id":          m.ID,
			"title":       m.Title,
			"slug":        m.Slug,
			"description": m.Description,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// ========================
// GET /api/materi/:slug
// ========================
func GetMateriDetail(c *gin.Context) {
	slug := c.Param("slug")

	materi, sections, err := GetMateriBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "materi not found",
		})
		return
	}

	// ambil user_id dari OptionalJWT (jika ada)
	var userID *uint = nil
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}

	isProUser := IsUserPro(userID)

	respSections := []gin.H{}

	for _, s := range sections {
		// SECTION FREE → selalu kirim content
		if s.IsFree {
			respSections = append(respSections, gin.H{
				"id":      s.ID,
				"title":   s.Title,
				"content": s.Content,
				"locked":  false,
			})
			continue
		}

		// SECTION PRO
		if isProUser {
			respSections = append(respSections, gin.H{
				"id":      s.ID,
				"title":   s.Title,
				"content": s.Content,
				"locked":  false,
			})
		} else {
			// konten TIDAK DIKIRIM
			respSections = append(respSections, gin.H{
				"id":     s.ID,
				"title":  s.Title,
				"locked": true,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"materi": gin.H{
			"id":          materi.ID,
			"title":       materi.Title,
			"slug":        materi.Slug,
			"description": materi.Description,
		},
		"sections": respSections,
	})
}

func UploadMateriImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image file is required",
		})
		return
	}

	// max 3MB
	if file.Size > 3*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image too large (max 3MB)",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid image type",
		})
		return
	}

	filename := fmt.Sprintf(
		"materi_%d%s",
		time.Now().UnixNano(),
		ext,
	)

	savePath := filepath.Join("public/materi/images", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save image",
		})
		return
	}

	imageURL := "/public/materi/images/" + filename

	c.JSON(http.StatusOK, gin.H{
		"image_url": imageURL,
	})
}

type CreateMateriRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
}

func AdminCreateMateri(c *gin.Context) {
	var req CreateMateriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	materi := models.Materi{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := CreateMateri(&materi); err != nil {
		c.JSON(500, gin.H{"error": "failed to create materi"})
		return
	}

	c.JSON(201, materi)
}

// =====================
// ADMIN UPDATE MATERI
// =====================
type UpdateMateriRequest struct {
	Title       *string `json:"title"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
}

func AdminUpdateMateri(c *gin.Context) {
	materiID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req UpdateMateriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if len(updates) == 0 {
		c.JSON(400, gin.H{"error": "no fields to update"})
		return
	}

	if err := UpdateMateri(uint(materiID), updates); err != nil {
		c.JSON(500, gin.H{"error": "failed to update materi"})
		return
	}

	c.JSON(200, gin.H{"message": "materi updated"})
}

// =====================
// ADMIN DELETE MATERI
// =====================
func AdminDeleteMateri(c *gin.Context) {
	materiID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := DeleteMateri(uint(materiID)); err != nil {
		c.JSON(500, gin.H{"error": "failed to delete materi"})
		return
	}

	c.JSON(200, gin.H{"message": "materi deleted"})
}

// =====================
// GET sections by materi
// =====================
func AdminGetSections(c *gin.Context) {
	materiID, _ := strconv.ParseUint(c.Param("materiId"), 10, 64)

	sections, err := GetSectionsByMateriID(uint(materiID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch sections",
		})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// =====================
// CREATE section
// =====================
type CreateSectionRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Order   int    `json:"order"`
	IsFree  bool   `json:"is_free"`
}

func AdminCreateSection(c *gin.Context) {
	materiID, _ := strconv.ParseUint(c.Param("materiId"), 10, 64)

	var req CreateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section := models.MateriSection{
		MateriID: uint(materiID),
		Title:    req.Title,
		Content:  req.Content,
		Order:    req.Order,
		IsFree:   req.IsFree,
	}

	if err := CreateSection(&section); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create section",
		})
		return
	}

	c.JSON(http.StatusCreated, section)
}

// =====================
// UPDATE section
// =====================
type UpdateSectionRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Order   *int    `json:"order"`
	IsFree  *bool   `json:"is_free"`
}

func AdminUpdateSection(c *gin.Context) {
	sectionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Order != nil {
		updates["order"] = *req.Order
	}
	if req.IsFree != nil {
		updates["is_free"] = *req.IsFree
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fields to update",
		})
		return
	}

	if err := UpdateSection(uint(sectionID), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update section",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "section updated",
	})
}

// =====================
// DELETE section
// =====================
func AdminDeleteSection(c *gin.Context) {
	sectionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := DeleteSection(uint(sectionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete section",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "section deleted",
	})
}
