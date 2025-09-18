package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type SimpanPinjamHandler struct {
	simpanPinjamService *services.SimpanPinjamService
}

func NewSimpanPinjamHandler(simpanPinjamService *services.SimpanPinjamService) *SimpanPinjamHandler {
	return &SimpanPinjamHandler{simpanPinjamService: simpanPinjamService}
}

func (h *SimpanPinjamHandler) CreateProduk(c *gin.Context) {
	var req services.CreateProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	produk, err := h.simpanPinjamService.CreateProduk(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk created successfully",
		"produk":  produk,
	})
}

func (h *SimpanPinjamHandler) GetProdukList(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	jenis := c.Query("jenis")

	produks, err := h.simpanPinjamService.GetProdukList(koperasiID, jenis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"produks": produks,
	})
}

func (h *SimpanPinjamHandler) CreateRekening(c *gin.Context) {
	var req services.CreateRekeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rekening, err := h.simpanPinjamService.CreateRekening(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Rekening created successfully",
		"rekening": rekening,
	})
}

func (h *SimpanPinjamHandler) GetRekeningByAnggota(c *gin.Context) {
	anggotaIDStr := c.Param("anggota_id")
	anggotaID, err := strconv.ParseUint(anggotaIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anggota ID"})
		return
	}

	rekenings, err := h.simpanPinjamService.GetRekeningByAnggota(anggotaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rekenings": rekenings,
	})
}

func (h *SimpanPinjamHandler) CreateTransaksi(c *gin.Context) {
	var req services.CreateTransaksiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if exists {
		if uid, ok := userID.(uint64); ok {
			req.CreatedBy = uid
		}
	}

	transaksi, err := h.simpanPinjamService.CreateTransaksi(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Transaksi created successfully",
		"transaksi": transaksi,
	})
}

func (h *SimpanPinjamHandler) GetTransaksiByRekening(c *gin.Context) {
	rekeningIDStr := c.Param("rekening_id")
	rekeningID, err := strconv.ParseUint(rekeningIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rekening ID"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	transaksis, err := h.simpanPinjamService.GetTransaksiByRekening(rekeningID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaksis": transaksis,
		"page":       page,
		"limit":      limit,
	})
}

func (h *SimpanPinjamHandler) GetStatistik(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	statistik, err := h.simpanPinjamService.GetStatistik(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statistik": statistik,
	})
}

func (h *SimpanPinjamHandler) GetPinjamanJatuhTempo(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 0 {
		days = 7
	}

	rekenings, err := h.simpanPinjamService.GetPinjamanJatuhTempo(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rekenings": rekenings,
		"days":      days,
	})
}