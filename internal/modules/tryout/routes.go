package tryout

import (
	"pintuniv-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// ADMIN
		api.POST("/admin/tryouts", AdminCreateTryout)
		api.POST("/admin/tryouts/:id/questions", AdminCreateTryoutQuestion)
		api.POST("/admin/tryout-questions/:id/options", AdminCreateTryoutOption)

		// USER
		api.POST(
			"/tryouts/:id/start",
			middleware.JWTAuth(),
			StartTryout,
		)
		api.GET(
			"/tryouts/:id/question/:number",
			middleware.JWTAuth(),
			GetTryoutQuestionByNumber,
		)
		api.POST(
			"/tryouts/:id/submit",
			middleware.JWTAuth(),
			SubmitTryout,
		)
		api.GET(
			"/tryouts/:id/resume",
			middleware.JWTAuth(),
			ResumeTryout,
		)

		api.POST(
			"/tryouts/:id/answer",
			middleware.JWTAuth(),
			SaveTryoutAnswer,
		)

	}
}
