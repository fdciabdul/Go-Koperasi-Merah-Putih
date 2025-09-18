package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type FinancialHandler struct {
	financialService *services.FinancialService
}

func NewFinancialHandler(financialService *services.FinancialService) *FinancialHandler {
	return &FinancialHandler{financialService: financialService}
}

func (h *FinancialHandler) CreateCOAAkun(c *gin.Context) {
	var req services.CreateCOAAkunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	akun, err := h.financialService.CreateCOAAkun(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "COA akun created successfully",
		"akun":    akun,
	})
}

func (h *FinancialHandler) GetCOAAkunList(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	akuns, err := h.financialService.GetCOAAkunList(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"akuns": akuns,
	})
}

func (h *FinancialHandler) GetCOAKategoriList(c *gin.Context) {
	kategoris, err := h.financialService.GetCOAKategoriList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kategoris": kategoris,
	})
}

func (h *FinancialHandler) CreateJurnal(c *gin.Context) {
	var req services.CreateJurnalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jurnal, err := h.financialService.CreateJurnalUmum(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Jurnal created successfully",
		"jurnal":  jurnal,
	})
}

func (h *FinancialHandler) GetJurnalList(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	dariStr := c.Query("dari")
	sampaiStr := c.Query("sampai")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	dari, err := time.Parse("2006-01-02", dariStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dari date format"})
		return
	}

	sampai, err := time.Parse("2006-01-02", sampaiStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sampai date format"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	jurnals, err := h.financialService.GetJurnalUmumList(koperasiID, dari, sampai, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jurnals": jurnals,
		"page":    page,
		"limit":   limit,
	})
}

func (h *FinancialHandler) GetJurnal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jurnal ID"})
		return
	}

	jurnal, err := h.financialService.GetJurnalUmumByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jurnal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jurnal": jurnal,
	})
}

func (h *FinancialHandler) PostJurnal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jurnal ID"})
		return
	}

	userID, _ := c.Get("user_id")

	err = h.financialService.PostJurnal(id, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jurnal posted successfully",
	})
}

func (h *FinancialHandler) CancelJurnal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jurnal ID"})
		return
	}

	userID, _ := c.Get("user_id")

	err = h.financialService.CancelJurnal(id, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jurnal cancelled successfully",
	})
}

func (h *FinancialHandler) GetNeracaSaldo(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	tanggalStr := c.Query("tanggal")
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tanggal format"})
		return
	}

	neracaSaldo, err := h.financialService.GetNeracaSaldo(koperasiID, tanggal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"neraca_saldo": neracaSaldo,
		"tanggal":      tanggal.Format("2006-01-02"),
	})
}

func (h *FinancialHandler) GetLabaRugi(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	dariStr := c.Query("dari")
	sampaiStr := c.Query("sampai")

	dari, err := time.Parse("2006-01-02", dariStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dari date format"})
		return
	}

	sampai, err := time.Parse("2006-01-02", sampaiStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sampai date format"})
		return
	}

	labaRugi, err := h.financialService.GetLabaRugi(koperasiID, dari, sampai)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"laba_rugi": labaRugi,
		"periode": gin.H{
			"dari":   dari.Format("2006-01-02"),
			"sampai": sampai.Format("2006-01-02"),
		},
	})
}

func (h *FinancialHandler) GetNeraca(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	tanggalStr := c.Query("tanggal")
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tanggal format"})
		return
	}

	neraca, err := h.financialService.GetNeraca(koperasiID, tanggal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"neraca":  neraca,
		"tanggal": tanggal.Format("2006-01-02"),
	})
}

func (h *FinancialHandler) GetSaldoAkun(c *gin.Context) {
	akunIDStr := c.Param("akun_id")
	akunID, err := strconv.ParseUint(akunIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid akun ID"})
		return
	}

	tanggalStr := c.Query("tanggal")
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tanggal format"})
		return
	}

	saldo, err := h.financialService.GetSaldoAkun(akunID, tanggal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"akun_id": akunID,
		"saldo":   saldo,
		"tanggal": tanggal.Format("2006-01-02"),
	})
}