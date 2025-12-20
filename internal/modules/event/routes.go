package event

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// ADMIN
		api.POST("/admin/events", AdminCreateEvent)
		api.POST("/admin/events/:id/tryouts", AdminAttachTryout)

		// USER
		api.GET("/events", ListActiveEvents)
		api.GET("/events/:id/tryouts", ListEventTryouts)
	}
}
