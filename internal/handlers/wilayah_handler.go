package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type WilayahHandler struct {
	wilayahService *services.WilayahService
}

func NewWilayahHandler(wilayahService *services.WilayahService) *WilayahHandler {
	return &WilayahHandler{wilayahService: wilayahService}
}

func (h *WilayahHandler) GetProvinsiList(c *gin.Context) {
	provinsis, err := h.wilayahService.GetProvinsiList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"provinsis": provinsis,
	})
}

func (h *WilayahHandler) GetKabupatenList(c *gin.Context) {
	provinsiIDStr := c.Param("provinsi_id")
	provinsiID, err := strconv.ParseUint(provinsiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provinsi ID"})
		return
	}

	kabupatens, err := h.wilayahService.GetKabupatenByProvinsi(provinsiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kabupatens": kabupatens,
	})
}

func (h *WilayahHandler) GetKecamatanList(c *gin.Context) {
	kabupatenIDStr := c.Param("kabupaten_id")
	kabupatenID, err := strconv.ParseUint(kabupatenIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kabupaten ID"})
		return
	}

	kecamatans, err := h.wilayahService.GetKecamatanByKabupaten(kabupatenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kecamatans": kecamatans,
	})
}

func (h *WilayahHandler) GetKelurahanList(c *gin.Context) {
	kecamatanIDStr := c.Param("kecamatan_id")
	kecamatanID, err := strconv.ParseUint(kecamatanIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kecamatan ID"})
		return
	}

	kelurahans, err := h.wilayahService.GetKelurahanByKecamatan(kecamatanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kelurahans": kelurahans,
	})
}