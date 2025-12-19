package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/utils"
)

func Register(name, email, password, phone string) error {
	// cek email sudah ada
	var existing models.User
	if err := database.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("email already registered")
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: string(hashed),
		Profile: models.Profile{
			FullName: name,
			Phone:    phone,
			IsPro:    false,
		},
	}

	// ====== TRANSACTION START ======
	tx := database.DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := GenerateOTP(user); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
	// ====== TRANSACTION END ======
}

func Login(email, password string) (string, string, error) {
	var user models.User

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	if !user.IsVerified {
		return "", "", errors.New("otp_required")
	}

	access, _ := utils.GenerateAccessToken(user.ID, user.Email)
	refresh, _ := utils.GenerateRefreshToken(user.ID)

	// simpan refresh token ke DB
	database.DB.Model(&user).Update("refresh_token", refresh)

	return access, refresh, nil
}

func Logout(userID uint) error {
	return database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("refresh_token", "").
		Error
}
