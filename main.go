package main

import (
	"os"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/modules/auth"
	"pintuniv-go/internal/modules/materi"
	"pintuniv-go/internal/modules/profile"

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
		&models.Materi{},
		&models.MateriSection{},
	)

	r := gin.Default()
	r.Static("/public", "./public")

	// register routes
	auth.RegisterRoutes(r)
	profile.RegisterRoutes(r)
	materi.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	r.Run(":" + port)
}
