package latihan

import (
	"pintuniv-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// USER
		api.GET("/materi/id/:id/latihan", ListLatihanByMateri)
		api.GET("/latihan/:id", middleware.OptionalJWT(), GetLatihanDetail)

		api.GET(
			"/latihan/:id/question/:number",
			middleware.JWTAuth(),
			GetQuestionByNumber,
		)

		// ADMIN
		api.POST("/admin/materi/id/:id/latihan", AdminCreateLatihan)
		api.POST("/admin/latihan/:id/questions", AdminCreateQuestion)
		api.POST("/admin/questions/:id/options", AdminCreateOption)

		api.POST(
			"/latihan/:id/start",
			middleware.JWTAuth(),
			StartLatihan,
		)
		api.POST(
			"/questions/:id/answer",
			middleware.JWTAuth(),
			SubmitAnswer,
		)
		api.POST(
			"/latihan/attempt/:attempt_id/finish",
			middleware.JWTAuth(),
			FinishLatihan,
		)

	}
}
