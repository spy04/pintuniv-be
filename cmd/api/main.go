// package main

// import (
// 	"os"

// 	"pintuniv-go/internal/config"
// 	"pintuniv-go/internal/database"
// 	"pintuniv-go/internal/modules/auth"
// 	"pintuniv-go/internal/modules/countdown"
// 	"pintuniv-go/internal/modules/event"
// 	"pintuniv-go/internal/modules/latihan"
// 	"pintuniv-go/internal/modules/materi"
// 	package_module "pintuniv-go/internal/modules/package"
// 	"pintuniv-go/internal/modules/payment"
// 	"pintuniv-go/internal/modules/profile"
// 	"pintuniv-go/internal/modules/promo"
// 	"pintuniv-go/internal/modules/quote"
// 	"pintuniv-go/internal/modules/tryout"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	config.InitMidtrans()
// 	database.ConnectDB()

// 	gin.SetMode(gin.ReleaseMode)
// 	r := gin.New()
// 	r.Use(gin.Recovery())

// 	r.Static("/public", "public")

// 	auth.RegisterRoutes(r)
// 	profile.RegisterRoutes(r)
// 	materi.RegisterRoutes(r)
// 	latihan.RegisterRoutes(r)
// 	tryout.RegisterRoutes(r)
// 	event.RegisterRoutes(r)
// 	payment.RegisterRoutes(r)
// 	package_module.RegisterRoutes(r)
// 	promo.RegisterRoutes(r)
// 	quote.RegisterRoutes(r)
// 	countdown.RegisterRoutes(r)
// 	// dst...

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}
// 	r.Run(":" + port)
// }

package main

import (
	"log"
	"os"

	"pintuniv-go/internal/config"
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/modules/auth"
	"pintuniv-go/internal/modules/countdown"
	"pintuniv-go/internal/modules/event"
	"pintuniv-go/internal/modules/latihan"
	"pintuniv-go/internal/modules/materi"
	package_module "pintuniv-go/internal/modules/package"
	"pintuniv-go/internal/modules/payment"
	"pintuniv-go/internal/modules/profile"
	"pintuniv-go/internal/modules/promo"
	"pintuniv-go/internal/modules/quote"
	"pintuniv-go/internal/modules/tryout"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// ✅ load .env hanya di local / dev
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ no .env file found, using system env")
		}
	}

	config.InitMidtrans()
	database.ConnectDB()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.Static("/public", "public")

	auth.RegisterRoutes(r)
	profile.RegisterRoutes(r)
	materi.RegisterRoutes(r)
	latihan.RegisterRoutes(r)
	tryout.RegisterRoutes(r)
	event.RegisterRoutes(r)
	payment.RegisterRoutes(r)
	package_module.RegisterRoutes(r)
	promo.RegisterRoutes(r)
	quote.RegisterRoutes(r)
	countdown.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
