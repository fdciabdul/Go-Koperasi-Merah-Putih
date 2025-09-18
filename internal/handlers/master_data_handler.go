package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type MasterDataHandler struct {
	masterDataService *services.MasterDataService
}

func NewMasterDataHandler(masterDataService *services.MasterDataService) *MasterDataHandler {
	return &MasterDataHandler{masterDataService: masterDataService}
}

func (h *MasterDataHandler) CreateKBLI(c *gin.Context) {
	var req services.CreateKBLIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kbli, err := h.masterDataService.CreateKBLI(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "KBLI created successfully",
		"kbli":    kbli,
	})
}

func (h *MasterDataHandler) GetKBLIList(c *gin.Context) {
	search := c.Query("search")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 50
	}

	kblis, err := h.masterDataService.GetKBLIList(search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kblis": kblis,
		"page":  page,
		"limit": limit,
	})
}

func (h *MasterDataHandler) GetKBLI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KBLI ID"})
		return
	}

	kbli, err := h.masterDataService.GetKBLIByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "KBLI not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kbli": kbli,
	})
}

func (h *MasterDataHandler) UpdateKBLI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KBLI ID"})
		return
	}

	var req services.UpdateKBLIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kbli, err := h.masterDataService.UpdateKBLI(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KBLI updated successfully",
		"kbli":    kbli,
	})
}

func (h *MasterDataHandler) CreateJenisKoperasi(c *gin.Context) {
	var req services.CreateJenisKoperasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jenis, err := h.masterDataService.CreateJenisKoperasi(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Jenis koperasi created successfully",
		"jenis":   jenis,
	})
}

func (h *MasterDataHandler) GetJenisKoperasiList(c *gin.Context) {
	jenisKoperasis, err := h.masterDataService.GetJenisKoperasiList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jenis_koperasi": jenisKoperasis,
	})
}

func (h *MasterDataHandler) GetJenisKoperasi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jenis koperasi ID"})
		return
	}

	jenis, err := h.masterDataService.GetJenisKoperasiByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jenis koperasi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jenis": jenis,
	})
}

func (h *MasterDataHandler) UpdateJenisKoperasi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jenis koperasi ID"})
		return
	}

	var req services.UpdateJenisKoperasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jenis, err := h.masterDataService.UpdateJenisKoperasi(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jenis koperasi updated successfully",
		"jenis":   jenis,
	})
}

func (h *MasterDataHandler) CreateBentukKoperasi(c *gin.Context) {
	var req services.CreateBentukKoperasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bentuk, err := h.masterDataService.CreateBentukKoperasi(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bentuk koperasi created successfully",
		"bentuk":  bentuk,
	})
}

func (h *MasterDataHandler) GetBentukKoperasiList(c *gin.Context) {
	bentukKoperasis, err := h.masterDataService.GetBentukKoperasiList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bentuk_koperasi": bentukKoperasis,
	})
}

func (h *MasterDataHandler) GetBentukKoperasi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bentuk koperasi ID"})
		return
	}

	bentuk, err := h.masterDataService.GetBentukKoperasiByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bentuk koperasi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bentuk": bentuk,
	})
}

func (h *MasterDataHandler) UpdateBentukKoperasi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bentuk koperasi ID"})
		return
	}

	var req services.UpdateBentukKoperasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bentuk, err := h.masterDataService.UpdateBentukKoperasi(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bentuk koperasi updated successfully",
		"bentuk":  bentuk,
	})
}