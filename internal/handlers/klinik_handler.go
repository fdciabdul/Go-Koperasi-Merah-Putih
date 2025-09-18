package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type KlinikHandler struct {
	klinikService *services.KlinikService
}

func NewKlinikHandler(klinikService *services.KlinikService) *KlinikHandler {
	return &KlinikHandler{klinikService: klinikService}
}

func (h *KlinikHandler) CreatePasien(c *gin.Context) {
	var req services.CreatePasienRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pasien, err := h.klinikService.CreatePasien(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pasien created successfully",
		"pasien":  pasien,
	})
}

func (h *KlinikHandler) GetPasien(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pasien ID"})
		return
	}

	pasien, err := h.klinikService.GetPasienByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pasien not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pasien": pasien,
	})
}

func (h *KlinikHandler) GetPasienList(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
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

	pasiens, err := h.klinikService.GetPasienList(koperasiID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pasiens": pasiens,
		"page":    page,
		"limit":   limit,
	})
}

func (h *KlinikHandler) SearchPasien(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	search := c.Query("q")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	pasiens, err := h.klinikService.SearchPasien(koperasiID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pasiens": pasiens,
	})
}

func (h *KlinikHandler) CreateTenagaMedis(c *gin.Context) {
	var req services.CreateTenagaMedisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenagaMedis, err := h.klinikService.CreateTenagaMedis(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Tenaga medis created successfully",
		"tenaga_medis": tenagaMedis,
	})
}

func (h *KlinikHandler) GetTenagaMedisList(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	tenagaMedis, err := h.klinikService.GetTenagaMedisList(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenaga_medis": tenagaMedis,
	})
}

func (h *KlinikHandler) CreateKunjungan(c *gin.Context) {
	var req services.CreateKunjunganRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kunjungan, err := h.klinikService.CreateKunjungan(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Kunjungan created successfully",
		"kunjungan": kunjungan,
	})
}

func (h *KlinikHandler) GetKunjungan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kunjungan ID"})
		return
	}

	kunjungan, err := h.klinikService.GetKunjunganByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kunjungan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kunjungan": kunjungan,
	})
}

func (h *KlinikHandler) GetKunjunganByPasien(c *gin.Context) {
	pasienIDStr := c.Param("pasien_id")
	pasienID, err := strconv.ParseUint(pasienIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pasien ID"})
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

	kunjungans, err := h.klinikService.GetKunjunganByPasien(pasienID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kunjungans": kunjungans,
		"page":       page,
		"limit":      limit,
	})
}

func (h *KlinikHandler) CreateObat(c *gin.Context) {
	var req services.CreateObatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	obat, err := h.klinikService.CreateObat(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Obat created successfully",
		"obat":    obat,
	})
}

func (h *KlinikHandler) GetObatList(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	obats, err := h.klinikService.GetObatList(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"obats": obats,
	})
}

func (h *KlinikHandler) SearchObat(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	search := c.Query("q")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	obats, err := h.klinikService.SearchObat(koperasiID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"obats": obats,
	})
}

func (h *KlinikHandler) GetStatistik(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	statistik, err := h.klinikService.GetStatistik(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statistik": statistik,
	})
}

func (h *KlinikHandler) GetObatStokRendah(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	obats, err := h.klinikService.GetObatStokRendah(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"obats": obats,
	})
}