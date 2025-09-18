package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	userService    *services.UserService
	ppobService    *services.PPOBService
}

func NewPaymentHandler(
	paymentService *services.PaymentService,
	userService *services.UserService,
	ppobService *services.PPOBService,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		userService:    userService,
		ppobService:    ppobService,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req services.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentService.CreatePayment(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Payment created successfully",
		"payment": payment,
	})
}

func (h *PaymentHandler) HandleMidtransCallback(c *gin.Context) {
	var callbackData map[string]interface{}
	if err := c.ShouldBindJSON(&callbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.paymentService.HandleCallback("midtrans", callbackData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.processPaymentCallback(callbackData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Callback processed successfully"})
}

func (h *PaymentHandler) HandleXenditCallback(c *gin.Context) {
	var callbackData map[string]interface{}
	if err := c.ShouldBindJSON(&callbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.paymentService.HandleCallback("xendit", callbackData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.processPaymentCallback(callbackData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Callback processed successfully"})
}

func (h *PaymentHandler) processPaymentCallback(callbackData map[string]interface{}) error {
	transactionID, _ := callbackData["order_id"].(string)
	if transactionID == "" {
		transactionID, _ = callbackData["external_id"].(string)
	}

	transactionStatus, _ := callbackData["transaction_status"].(string)
	if transactionStatus == "" {
		transactionStatus, _ = callbackData["status"].(string)
	}

	if transactionStatus == "settlement" || transactionStatus == "capture" || transactionStatus == "PAID" {
		return nil
	}

	return nil
}