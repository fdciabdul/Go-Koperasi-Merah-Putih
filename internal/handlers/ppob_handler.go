package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type PPOBHandler struct {
	ppobService *services.PPOBService
}

func NewPPOBHandler(ppobService *services.PPOBService) *PPOBHandler {
	return &PPOBHandler{ppobService: ppobService}
}

func (h *PPOBHandler) GetKategoriList(c *gin.Context) {
	kategoris, err := h.ppobService.GetKategoriList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Success",
		"kategoris": kategoris,
	})
}

func (h *PPOBHandler) GetProdukByKategori(c *gin.Context) {
	kategoriIDStr := c.Param("kategori_id")
	kategoriID, err := strconv.ParseUint(kategoriIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kategori ID"})
		return
	}

	produks, err := h.ppobService.GetProdukByKategori(kategoriID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"produks": produks,
	})
}

func (h *PPOBHandler) CreateTransaction(c *gin.Context) {
	var req services.PPOBTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaksi, err := h.ppobService.CreateTransaction(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Transaction created successfully",
		"transaksi": transaksi,
	})
}

func (h *PPOBHandler) CreateSettlement(c *gin.Context) {
	var req struct {
		KoperasiID uint64 `json:"koperasi_id" binding:"required"`
		Dari       string `json:"dari" binding:"required"`
		Sampai     string `json:"sampai" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dari, err := time.Parse("2006-01-02", req.Dari)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format for 'dari'. Use YYYY-MM-DD"})
		return
	}

	sampai, err := time.Parse("2006-01-02", req.Sampai)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format for 'sampai'. Use YYYY-MM-DD"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	processedBy, ok := userID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	settlement, err := h.ppobService.CreateSettlement(req.KoperasiID, dari, sampai, processedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Settlement created successfully",
		"settlement": settlement,
	})
}