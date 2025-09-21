package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type KoperasiHandler struct {
	koperasiService *services.KoperasiService
}

func NewKoperasiHandler(koperasiService *services.KoperasiService) *KoperasiHandler {
	return &KoperasiHandler{koperasiService: koperasiService}
}

func (h *KoperasiHandler) CreateKoperasi(c *gin.Context) {
	var req services.CreateKoperasiRequest
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

	koperasi, err := h.koperasiService.CreateKoperasi(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Koperasi created successfully",
		"koperasi": koperasi,
	})
}

func (h *KoperasiHandler) GetKoperasi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	koperasi, err := h.koperasiService.GetKoperasiByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Koperasi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"koperasi": koperasi,
	})
}

func (h *KoperasiHandler) GetKoperasiList(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID required"})
		return
	}

	tid, ok := tenantID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid tenant ID"})
		return
	}

	koperasis, err := h.koperasiService.GetKoperasiByTenant(tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"koperasis": koperasis,
	})
}

func (h *KoperasiHandler) UpdateKoperasi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	var req services.UpdateKoperasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if exists {
		if uid, ok := userID.(uint64); ok {
			req.UpdatedBy = uid
		}
	}

	koperasi, err := h.koperasiService.UpdateKoperasi(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Koperasi updated successfully",
		"koperasi": koperasi,
	})
}

func (h *KoperasiHandler) DeleteKoperasi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
		return
	}

	err = h.koperasiService.DeleteKoperasi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Koperasi deleted successfully",
	})
}

func (h *KoperasiHandler) CreateAnggota(c *gin.Context) {
	var req services.CreateAnggotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anggota, err := h.koperasiService.CreateAnggota(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Anggota created successfully",
		"anggota": anggota,
	})
}

func (h *KoperasiHandler) GetAnggota(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anggota ID"})
		return
	}

	anggota, err := h.koperasiService.GetAnggotaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anggota not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"anggota": anggota,
	})
}

func (h *KoperasiHandler) GetAnggotaList(c *gin.Context) {
	koperasiIDStr := c.Param("id")
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

	anggotas, err := h.koperasiService.GetAnggotaByKoperasi(koperasiID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"anggotas": anggotas,
		"page":     page,
		"limit":    limit,
	})
}

func (h *KoperasiHandler) UpdateAnggotaStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anggota ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=aktif non_aktif keluar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.koperasiService.UpdateAnggotaStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Anggota status updated successfully",
	})
}

func (h *KoperasiHandler) GetProvinsiList(c *gin.Context) {
	provinsis, err := h.koperasiService.GetProvinsiList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"provinsis": provinsis,
	})
}

func (h *KoperasiHandler) GetKabupatenList(c *gin.Context) {
	provinsiIDStr := c.Param("provinsi_id")
	provinsiID, err := strconv.ParseUint(provinsiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provinsi ID"})
		return
	}

	kabupatens, err := h.koperasiService.GetKabupatenByProvinsi(provinsiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kabupatens": kabupatens,
	})
}

func (h *KoperasiHandler) GetKecamatanList(c *gin.Context) {
	kabupatenIDStr := c.Param("kabupaten_id")
	kabupatenID, err := strconv.ParseUint(kabupatenIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kabupaten ID"})
		return
	}

	kecamatans, err := h.koperasiService.GetKecamatanByKabupaten(kabupatenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kecamatans": kecamatans,
	})
}

func (h *KoperasiHandler) GetKelurahanList(c *gin.Context) {
	kecamatanIDStr := c.Param("kecamatan_id")
	kecamatanID, err := strconv.ParseUint(kecamatanIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kecamatan ID"})
		return
	}

	kelurahans, err := h.koperasiService.GetKelurahanByKecamatan(kecamatanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kelurahans": kelurahans,
	})
}