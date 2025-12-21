// package main

// import (
// 	"pintuniv-go/internal/database"
// 	"pintuniv-go/internal/database/migrations"
// 	"pintuniv-go/internal/models"
// )

// func main() {
// 	database.ConnectDB()

// 	database.DB.AutoMigrate(
// 		&models.User{},
// 		&models.Profile{},
// 		&models.OTP{},
// 		&models.PasswordReset{},
// 		&models.Materi{},
// 		&models.MateriSection{},
// 		&models.Latihan{},
// 		&models.Option{},
// 		&models.Question{},
// 		&models.UserAnswer{},
// 		&models.UserLatihanAttempt{},
// 		&models.Tryout{},
// 		&models.TryoutQuestion{},
// 		&models.TryoutOption{},
// 		&models.TryoutAttempt{},
// 		&models.TryoutAnswer{},

// 		&models.Event{},
// 		&models.EventTryout{},

// 		&models.Payment{},
// 		&models.Package{},

// 		&models.Promo{},
// 		&models.PromoPackage{},

// 		&models.Quote{},

// 		&models.Countdown{},
// 	)

// 	migrations.Run()
// }

package main

import (
	"log"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/database/migrations"
	"pintuniv-go/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	// load env for local migration
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ no .env file found, using system env")
	}

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

		&models.Promo{},
		&models.PromoPackage{},

		&models.Quote{},

		&models.Countdown{},
	)

	migrations.Run()
}
