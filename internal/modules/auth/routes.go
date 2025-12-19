package auth

import (
	"pintuniv-go/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// rate limiters
	loginLimiter := middleware.NewRateLimiter(10, time.Minute)
	otpLimiter := middleware.NewRateLimiter(5, time.Minute)
	{
		api.POST("/register", RegisterHandler)
		api.POST("/login", loginLimiter.Middleware(), LoginHandler)

		api.GET("/me", middleware.JWTAuth(), MeHandler)
		api.POST("/refresh", RefreshHandler)
		api.POST("/logout", middleware.JWTAuth(), LogoutHandler)

		api.POST("/verify-otp", otpLimiter.Middleware(), VerifyOTPHandler)
		api.POST("/resend-otp", otpLimiter.Middleware(), ResendOTPHandler)

		api.POST("/forgot-password", ForgotPasswordHandler)
		api.POST("/reset-password", ResetPasswordHandler)

	}
}
