package package_module

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// USER
		api.GET("/packages", ListPackages)

		// ADMIN
		api.POST("/admin/packages", AdminCreatePackage)
		api.PUT("/admin/packages/:id", AdminUpdatePackage)
		api.DELETE("/admin/packages/:id", AdminDeletePackage)
	}
}
