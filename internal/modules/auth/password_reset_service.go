package auth

import (
	"os"

	"crypto/rand"
	"encoding/hex"
	"time"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordResetLink(user models.User) error {
	// ===== 1. INVALIDATE TOKEN LAMA =====
	if err := database.DB.Model(&models.PasswordReset{}).
		Where("user_id = ? AND used = false", user.ID).
		Update("used", true).Error; err != nil {
		return err
	}

	// ===== 2. GENERATE TOKEN BARU =====
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	token := hex.EncodeToString(b)

	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	reset := models.PasswordReset{
		UserID:    user.ID,
		TokenHash: string(hash),
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Used:      false,
	}

	if err := database.DB.Create(&reset).Error; err != nil {
		return err
	}

	// ===== 3. BUILD LINK DARI ENV =====
	frontendURL := os.Getenv("FRONTEND_URL")
	link := frontendURL + "/reset-password?token=" + token

	// ===== 4. SEND EMAIL =====
	body := `
		<h2>Reset Password Pintuniv</h2>
		<p>Klik link di bawah untuk mengganti password kamu:</p>
		<p><a href="` + link + `">Reset Password</a></p>
		<p>Link ini berlaku selama <b>15 menit</b>.</p>
	`

	return utils.SendEmail(
		user.Email,
		"Reset Password Pintuniv",
		body,
	)
}
