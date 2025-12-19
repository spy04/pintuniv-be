package main

import (
	"os"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/modules/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	database.ConnectDB()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.OTP{},
		&models.PasswordReset{},
	)

	r := gin.Default()

	// register routes
	auth.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	r.Run(":" + port)
}
