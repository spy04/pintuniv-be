package payment

import (
	"fmt"

	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
	"pintuniv-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type CreateProPaymentRequest struct {
	PackageID uint   `json:"package_id" binding:"required"`
	PromoCode string `json:"promo_code"`
}

func CreateProPayment(c *gin.Context) {
	userID := c.GetUint("user_id")

	// --------------------------
	// 1. PARSE REQUEST
	// --------------------------
	var req CreateProPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// --------------------------
	// 2. AMBIL PACKAGE (SUMBER HARGA)
	// --------------------------
	var pkg models.Package
	if err := database.DB.
		Where("id = ? AND is_active = true", req.PackageID).
		First(&pkg).Error; err != nil {

		c.JSON(404, gin.H{"error": "package_not_found"})
		return
	}

	// --------------------------
	// 3. HITUNG HARGA AKHIR
	// --------------------------
	finalAmount := pkg.Price

	if req.PromoCode != "" {
		var promo models.Promo

		// promo harus aktif & dalam periode
		err := database.DB.
			Where(
				"code = ? AND is_active = true AND start_at <= NOW() AND end_at >= NOW()",
				req.PromoCode,
			).
			First(&promo).Error

		if err == nil {
			// cek promo berlaku untuk package ini atau tidak
			var count int64
			database.DB.
				Model(&models.PromoPackage{}).
				Where("promo_id = ? AND package_id = ?", promo.ID, pkg.ID).
				Count(&count)

			if count > 0 {
				if promo.Type == "percent" {
					discount := (pkg.Price * promo.Value) / 100
					finalAmount = pkg.Price - discount
				}

				if promo.Type == "flat" {
					finalAmount = pkg.Price - promo.Value
				}

				if finalAmount < 0 {
					finalAmount = 0
				}
			}
		}
	}

	// --------------------------
	// 4. CREATE PAYMENT RECORD
	// --------------------------
	orderID := fmt.Sprintf("PRO-%s", uuid.NewString())

	payment := models.Payment{
		UserID:    userID,
		PackageID: pkg.ID,
		OrderID:   orderID,
		Amount:    finalAmount,
		Status:    "pending",
		Type:      "pro",
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed_create_payment"})
		return
	}

	// --------------------------
	// 5. CREATE MIDTRANS SNAP
	// --------------------------
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: finalAmount,
		},
	}

	var s snap.Client
	s.New(midtrans.ServerKey, midtrans.Environment)

	resp, err := s.CreateTransaction(snapReq)
	if err != nil {
		c.JSON(500, gin.H{"error": "midtrans_error"})
		return
	}

	// --------------------------
	// 6. RESPONSE
	// --------------------------
	c.JSON(200, gin.H{
		"snap_url": resp.RedirectURL,
		"amount":   finalAmount,
	})
}

func MidtransCallback(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid_payload"})
		return
	}

	orderID := payload["order_id"].(string)
	statusCode := payload["status_code"].(string)
	grossAmount := payload["gross_amount"].(string)
	signatureKey := payload["signature_key"].(string)
	transactionStatus := payload["transaction_status"].(string)

	// =========================
	// VERIFY SIGNATURE
	// =========================
	if !utils.VerifyMidtransSignature(
		orderID,
		statusCode,
		grossAmount,
		midtrans.ServerKey,
		signatureKey,
	) {
		c.JSON(403, gin.H{"error": "invalid_signature"})
		return
	}

	// =========================
	// TRANSACTION START
	// =========================
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(500, gin.H{"error": "tx_failed"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var payment models.Payment
	if err := tx.
		Where("order_id = ?", orderID).
		First(&payment).Error; err != nil {

		tx.Rollback()
		c.JSON(404, gin.H{"error": "payment_not_found"})
		return
	}

	// =========================
	// IDEMPOTENCY CHECK
	// =========================
	if payment.Status == "success" || payment.Status == "failed" {
		tx.Rollback()
		c.JSON(200, gin.H{"status": "already_processed"})
		return
	}

	// =========================
	// STATUS TRANSITION
	// =========================
	switch transactionStatus {
	case "capture", "settlement":
		payment.Status = "success"

		if payment.Type == "pro" {
			tx.
				Model(&models.User{}).
				Where("id = ?", payment.UserID).
				Update("is_pro", true)
		}

	case "expire", "cancel":
		payment.Status = "failed"

	default:
		// pending → biarin
		tx.Rollback()
		c.JSON(200, gin.H{"status": "pending"})
		return
	}

	if err := tx.Save(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "failed_update_payment"})
		return
	}

	// =========================
	// COMMIT
	// =========================
	tx.Commit()

	c.JSON(200, gin.H{"status": "ok"})
}

func GetUserPaymentHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var payments []models.Payment
	database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&payments)

	resp := []gin.H{}
	for _, p := range payments {
		resp = append(resp, gin.H{
			"order_id":   p.OrderID,
			"package_id": p.PackageID,
			"type":       p.Type,
			"amount":     p.Amount,
			"status":     p.Status,
			"created_at": p.CreatedAt,
		})
	}

	c.JSON(200, resp)
}

func GetAllPayments(c *gin.Context) {
	status := c.Query("status") // optional
	pType := c.Query("type")    // optional

	query := database.DB.Model(&models.Payment{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if pType != "" {
		query = query.Where("type = ?", pType)
	}

	var payments []models.Payment
	query.
		Order("created_at DESC").
		Find(&payments)

	resp := []gin.H{}
	for _, p := range payments {
		resp = append(resp, gin.H{
			"id":         p.ID,
			"user_id":    p.UserID,
			"order_id":   p.OrderID,
			"type":       p.Type,
			"amount":     p.Amount,
			"status":     p.Status,
			"created_at": p.CreatedAt,
		})
	}

	c.JSON(200, resp)
}

func AdminCleanupPayments(c *gin.Context) {
	if err := ExpireOldPendingPayments(); err != nil {
		c.JSON(500, gin.H{"error": "cleanup_failed"})
		return
	}
	c.JSON(200, gin.H{"status": "cleaned"})
}
