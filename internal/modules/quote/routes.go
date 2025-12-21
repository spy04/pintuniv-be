package quote

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/quote", GetQuote)
		api.POST("/admin/quote", AdminCreateQuote)
		api.PUT("/admin/quote/:id", AdminUpdateQuote)
		api.DELETE("/admin/quote/:id", AdminDeleteQuote)

	}
}
