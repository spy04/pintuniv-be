package profile

import (
	"pintuniv-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api", middleware.JWTAuth())
	{
		api.GET("/profile", GetProfileHandler)
		api.PATCH("/profile", UpdateProfileHandler)
		api.POST("/profile/avatar", UploadAvatarHandler)

	}
}
