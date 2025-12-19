package tryout

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// ADMIN
		api.POST("/admin/tryouts", AdminCreateTryout)
		api.POST("/admin/tryouts/:id/questions", AdminCreateTryoutQuestion)
		api.POST("/admin/tryout-questions/:id/options", AdminCreateTryoutOption)

		// USER (nanti)
	}
}
