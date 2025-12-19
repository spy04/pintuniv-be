package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func GenerateOTP(user models.User) error {
	// cek cooldown
	if err := CanResendOTP(user.ID); err != nil {
		return err
	}

	code := rand.Intn(900000) + 100000
	codeStr := strconv.Itoa(code)

	hash, _ := bcrypt.GenerateFromPassword([]byte(codeStr), bcrypt.DefaultCost)

	otp := models.OTP{
		UserID:    user.ID,
		CodeHash:  string(hash),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := database.DB.Create(&otp).Error; err != nil {
		return err
	}

	body := fmt.Sprintf(`
		<h2>Kode OTP Pintuniv</h2>
		<h1>%s</h1>
		<p>Berlaku 5 menit.</p>
	`, codeStr)

	return utils.SendEmail(
		user.Email,
		"Kode OTP Pintuniv",
		body,
	)
}

func VerifyOTP(userID uint, code string) error {
	var otp models.OTP

	err := database.DB.
		Where("user_id = ? AND used = false", userID).
		Order("created_at desc").
		First(&otp).Error
	if err != nil {
		return err
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("otp expired")
	}

	if bcrypt.CompareHashAndPassword([]byte(otp.CodeHash), []byte(code)) != nil {
		return errors.New("invalid otp")
	}

	// tandai otp terpakai
	database.DB.Model(&otp).Update("used", true)

	// aktifkan user
	database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_verified", true)

	return nil
}

func CanResendOTP(userID uint) error {
	var lastOTP models.OTP

	err := database.DB.
		Where("user_id = ? AND used = false", userID).
		Order("created_at desc").
		First(&lastOTP).Error

	if err != nil {
		// belum pernah ada OTP → boleh
		return nil
	}

	// cooldown 60 detik
	if time.Since(lastOTP.CreatedAt) < time.Minute {
		return errors.New("please wait before requesting a new otp")
	}

	return nil
}
