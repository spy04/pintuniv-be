package main

import (
	"os"

	"pintuniv-go/internal/config"
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/database/migrations"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/modules/auth"
	"pintuniv-go/internal/modules/event"
	"pintuniv-go/internal/modules/latihan"
	"pintuniv-go/internal/modules/materi"
	package_module "pintuniv-go/internal/modules/package"
	"pintuniv-go/internal/modules/payment"
	"pintuniv-go/internal/modules/profile"
	"pintuniv-go/internal/modules/tryout"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.InitMidtrans()

	database.ConnectDB()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.OTP{},
		&models.PasswordReset{},
		&models.Materi{},
		&models.MateriSection{},
		&models.Latihan{},
		&models.Option{},
		&models.Question{},
		&models.UserAnswer{},
		&models.UserLatihanAttempt{},
		&models.Tryout{},
		&models.TryoutQuestion{},
		&models.TryoutOption{},
		&models.TryoutAttempt{},
		&models.TryoutAnswer{},

		&models.Event{},
		&models.EventTryout{},

		&models.Payment{},
		&models.Package{},
	)

	migrations.Run()

	r := gin.Default()
	r.Static("/public", "./public")

	// register routes
	auth.RegisterRoutes(r)
	profile.RegisterRoutes(r)
	materi.RegisterRoutes(r)
	latihan.RegisterRoutes(r)
	tryout.RegisterRoutes(r)
	event.RegisterRoutes(r)
	payment.RegisterRoutes(r)
	package_module.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	r.Run(":" + port)
}
