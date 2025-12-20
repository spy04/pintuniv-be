package payment

import (
	"pintuniv-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET(
			"/payments/history",
			middleware.JWTAuth(),
			GetUserPaymentHistory,
		)

		// ADMIN (sementara tanpa auth)
		api.GET(
			"/admin/payments",
			GetAllPayments,
		)
		api.POST(
			"/payments/pro",
			middleware.JWTAuth(),
			CreateProPayment,
		)

		api.POST(
			"/payments/midtrans/callback",
			MidtransCallback,
		)
		api.POST("/admin/payments/cleanup", AdminCleanupPayments)

	}
}
