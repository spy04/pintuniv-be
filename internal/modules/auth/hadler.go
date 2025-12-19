package auth

import (
	"net/http"
	"os"
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"` // optional
}

func RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := Register(req.Name, req.Email, req.Password, req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access":  access,
		"refresh": refresh,
	})
}

func MeHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	err := database.DB.Preload("Profile").
		First(&user, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"profile": gin.H{
			"name":   user.Profile.FullName,
			"phone":  user.Profile.Phone,
			"is_pro": user.Profile.IsPro,
		},
	})
}

type RefreshRequest struct {
	Refresh string `json:"refresh" binding:"required"`
}

func RefreshHandler(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(req.Refresh, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["type"] != "refresh" {
		c.JSON(401, gin.H{"error": "invalid token type"})
		return
	}
	userID := uint(claims["user_id"].(float64))

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// cocokkan refresh token
	if user.RefreshToken != req.Refresh {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token mismatch"})
		return
	}

	access, _ := utils.GenerateAccessToken(user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"access": access,
	})
}

func LogoutHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := Logout(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

func VerifyOTPHandler(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	if err := VerifyOTP(user.ID, req.OTP); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "otp verified"})
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func ResendOTPHandler(c *gin.Context) {
	var req ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	if user.IsVerified {
		c.JSON(400, gin.H{"error": "user already verified"})
		return
	}

	if err := GenerateOTP(user); err != nil {
		c.JSON(429, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "otp resent"})
}

func ForgotPasswordHandler(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err == nil {
		_ = GeneratePasswordResetLink(user)
	}

	// RESPONSE SELALU SAMA (anti email enumeration)
	c.JSON(200, gin.H{
		"message": "if email exists, reset link has been sent",
	})
}

func ResetPasswordHandler(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var resets []models.PasswordReset
	database.DB.
		Where("used = false").
		Find(&resets)

	var reset models.PasswordReset
	found := false
	for _, r := range resets {
		if bcrypt.CompareHashAndPassword(
			[]byte(r.TokenHash),
			[]byte(req.Token),
		) == nil {
			reset = r
			found = true
			break
		}
	}

	if !found {
		c.JSON(400, gin.H{"error": "invalid token"})
		return
	}

	if time.Now().After(reset.ExpiresAt) {
		c.JSON(400, gin.H{"error": "token expired"})
		return
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(reset.TokenHash),
		[]byte(req.Token),
	) != nil {
		c.JSON(400, gin.H{"error": "invalid token"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	tx := database.DB.Begin()
	tx.Model(&models.User{}).
		Where("id = ?", reset.UserID).
		Updates(map[string]interface{}{
			"password":      string(hashed),
			"refresh_token": "",
		})

	tx.Model(&reset).Update("used", true)
	tx.Commit()

	c.JSON(200, gin.H{"message": "password reset successful"})
}
