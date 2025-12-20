package profile

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
	IsPro    *bool   `json:"is_pro"`
}

func GetProfileHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	profile, err := GetProfileByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func UpdateProfileHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := UpdateProfile(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

func UploadAvatarHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 1. ambil avatar lama
	oldAvatarURL, _ := GetAvatarByUserID(userID)

	// 2. ambil file baru
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, gin.H{"error": "avatar file required"})
		return
	}

	if file.Size > 2*1024*1024 {
		c.JSON(400, gin.H{"error": "file too large (max 2MB)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(400, gin.H{"error": "invalid file type"})
		return
	}

	// 3. simpan avatar baru
	filename := fmt.Sprintf(
		"user_%d_%d%s",
		userID,
		time.Now().Unix(),
		ext,
	)

	savePath := filepath.Join("public/avatars", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(500, gin.H{"error": "failed to save avatar"})
		return
	}

	newAvatarURL := "/public/avatars/" + filename

	// 4. update DB
	if err := UpdateAvatar(userID, newAvatarURL); err != nil {
		c.JSON(500, gin.H{"error": "failed to update avatar"})
		return
	}

	// 5. hapus avatar lama (jika ada)
	if oldAvatarURL != "" {
		oldPath := "." + oldAvatarURL // "/public/avatars/xxx.jpg"
		_ = os.Remove(oldPath)        // ignore error (aman)
	}

	c.JSON(200, gin.H{
		"message":    "avatar updated",
		"avatar_url": newAvatarURL,
	})
}
