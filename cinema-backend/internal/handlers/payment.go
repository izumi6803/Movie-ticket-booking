package handlers

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	vnpayService   *services.VNPayService
	bookingService *services.BookingService
	paymentService *services.PaymentService
}

func NewPaymentHandler(vnpayService *services.VNPayService, bookingService *services.BookingService, paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		vnpayService:   vnpayService,
		bookingService: bookingService,
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreateVNPayPayment(c *gin.Context) {
	var request struct {
		BookingID string  `json:"bookingId"`
		Amount    float64 `json:"amount"`
		OrderInfo string  `json:"orderInfo"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "user not found in context"})
		return
	}

	// Convert userID to string
	userIDStr := fmt.Sprintf("%v", userID)

	// Create payment record
	payment, err := h.paymentService.CreatePayment(request.BookingID, userIDStr, request.Amount, models.PaymentMethodVNPay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Create VNPay payment URL
	response, err := h.vnpayService.CreatePayment(request.BookingID, request.Amount, request.OrderInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Update payment with order ID
	h.paymentService.UpdateOrderID(payment.ID.String(), response.OrderId)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"paymentUrl": response.PaymentUrl,
			"orderId":    response.OrderId,
			"paymentId":  payment.ID,
		},
	})
}

func (h *PaymentHandler) VNPayReturn(c *gin.Context) {
	// Get all query params
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	// Verify signature
	if !h.vnpayService.VerifyReturn(params) {
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/payment/callback?status=failed&message=invalid_signature")
		return
	}

	bookingID := params["vnp_OrderInfo"]
	responseCode := params["vnp_ResponseCode"]
	transactionID := params["vnp_TransactionNo"]
	orderID := params["vnp_TxnRef"]
	bankCode := params["vnp_BankCode"]

	if responseCode == "00" {
		// Payment successful
		paymentInfo := map[string]string{
			"transactionId": transactionID,
			"orderId":       orderID,
			"bankCode":      bankCode,
			"responseCode":  responseCode,
		}

		if err := h.bookingService.ConfirmPayment(bookingID, paymentInfo); err != nil {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://localhost:3000/payment/callback?status=failed&bookingId=%s&message=%s", bookingID, err.Error()))
			return
		}

		// Update payment record
		h.paymentService.UpdatePaymentStatusByOrderID(orderID, models.PaymentPaid, transactionID, bankCode)

		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://localhost:3000/payment/callback?status=success&bookingId=%s", bookingID))
	} else {
		// Payment failed
		h.bookingService.CancelBooking(bookingID, "payment_failed")
		h.paymentService.UpdatePaymentStatusByOrderID(orderID, models.PaymentFailed, "", "")

		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://localhost:3000/payment/callback?status=failed&bookingId=%s&code=%s", bookingID, responseCode))
	}
}

func (h *PaymentHandler) CheckPaymentStatus(c *gin.Context) {
	bookingID := c.Param("bookingId")

	booking, err := h.bookingService.GetByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"bookingId":     booking.ID,
			"status":        booking.Status,
			"paymentStatus": booking.PaymentStatus,
			"amount":        booking.TotalAmount,
			"bookingCode":   booking.BookingCode,
			"qrCode":        booking.QRCode,
		},
	})
}

// MockPaymentPage serves a simple HTML page for testing payments in sandbox mode
func (h *PaymentHandler) MockPaymentPage(c *gin.Context) {
	bookingID := c.Query("bookingId")
	amount := c.Query("amount")
	orderId := c.Query("orderId")

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>VNPay Test Payment</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; background: #f5f5f5; }
        .container { background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); text-align: center; max-width: 400px; }
        .logo { width: 80px; height: 80px; background: #e50019; border-radius: 20px; display: flex; align-items: center; justify-content: center; margin: 0 auto 20px; color: white; font-weight: bold; font-size: 20px; }
        .amount { font-size: 32px; font-weight: bold; color: #333; margin: 20px 0; }
        .info { color: #666; margin-bottom: 30px; }
        .btn { display: block; width: 100%%; padding: 15px; margin: 10px 0; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; text-decoration: none; }
        .btn-success { background: #4CAF50; color: white; }
        .btn-fail { background: #f44336; color: white; }
        .btn:hover { opacity: 0.9; }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">VNPay</div>
        <h2>Test Payment</h2>
        <div class="amount">%s VND</div>
        <div class="info">
            <p>Booking ID: %s</p>
            <p>Order ID: %s</p>
            <p>This is a sandbox payment for testing</p>
        </div>
        <a href="http://localhost:3001/api/payments/vnpay/mock-success?bookingId=%s&orderId=%s" class="btn btn-success">Simulate Successful Payment</a>
        <a href="http://localhost:3001/api/payments/vnpay/mock-fail?bookingId=%s&orderId=%s" class="btn btn-fail">Simulate Failed Payment</a>
    </div>
</body>
</html>`, amount, bookingID, orderId, bookingID, orderId, bookingID, orderId)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// MockPaymentSuccess simulates a successful payment
func (h *PaymentHandler) MockPaymentSuccess(c *gin.Context) {
	bookingID := c.Query("bookingId")
	orderID := c.Query("orderId")

	paymentInfo := map[string]string{
		"transactionId": "MOCK_" + orderID,
		"orderId":       orderID,
		"bankCode":      "NCB",
		"responseCode":  "00",
	}

	if err := h.bookingService.ConfirmPayment(bookingID, paymentInfo); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://localhost:3000/payment/callback?status=failed&bookingId=%s&message=%s", bookingID, err.Error()))
		return
	}

	// Update payment record
	h.paymentService.UpdatePaymentStatusByOrderID(orderID, models.PaymentPaid, "MOCK_"+orderID, "NCB")

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://localhost:3000/payment/callback?status=success&bookingId=%s", bookingID))
}

// MockPaymentFail simulates a failed payment
func (h *PaymentHandler) MockPaymentFail(c *gin.Context) {
	bookingID := c.Query("bookingId")
	orderID := c.Query("orderId")

	h.bookingService.CancelBooking(bookingID, "mock_payment_failed")
	h.paymentService.UpdatePaymentStatusByOrderID(orderID, models.PaymentFailed, "", "")

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://localhost:3000/payment/callback?status=failed&bookingId=%s&code=99", bookingID))
}
