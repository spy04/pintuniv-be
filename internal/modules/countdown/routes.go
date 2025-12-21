package countdown

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/admin/countdown", AdminUpsertCountdown)
		api.GET("/countdown", GetCountdown)
	}
}
