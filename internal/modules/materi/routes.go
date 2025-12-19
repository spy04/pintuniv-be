package materi

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/materi", GetMateriList)
		api.GET("/materi/:slug", GetMateriDetail)

		// admin materi
		api.POST("/admin/materi", AdminCreateMateri)
		api.PATCH("/admin/materi/:id", AdminUpdateMateri)
		api.DELETE("/admin/materi/:id", AdminDeleteMateri)

		// admin
		api.GET("/admin/materi/:materiId/sections", AdminGetSections)
		api.POST("/admin/materi/:materiId/sections", AdminCreateSection)
		api.PATCH("/admin/sections/:id", AdminUpdateSection)
		api.DELETE("/admin/sections/:id", AdminDeleteSection)

		api.POST("/admin/materi/upload-image", UploadMateriImage)

	}
}
