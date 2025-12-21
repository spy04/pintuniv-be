package promo

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// ======================
		// ADMIN PROMO
		// ======================

		// 1. Buat promo
		api.POST("/admin/promos", AdminCreatePromo)

		// 2. Attach promo ke package
		api.POST("/admin/promos/:id/packages", AdminAttachPromoPackage)

		// 3. Update promo (opsional)
		api.PUT("/admin/promos/:id", AdminUpdatePromo)

		// 4. Disable promo
		api.DELETE("/admin/promos/:id", AdminDisablePromo)

		// ======================
		// USER (OPSIONAL)
		// ======================

		// Cek promo valid / tidak (untuk FE preview)
		api.GET("/promos/validate", ValidatePromo)

		api.POST("/admin/promos/:id/image", AdminUploadPromoImage)

		api.GET("/promos", GetPromos)

	}
}
