package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/services"
)

type SequenceHandler struct {
	sequenceService *services.SequenceService
}

func NewSequenceHandler(sequenceService *services.SequenceService) *SequenceHandler {
	return &SequenceHandler{sequenceService: sequenceService}
}

func (h *SequenceHandler) GetSequenceList(c *gin.Context) {
	tenantIDInterface, _ := c.Get("tenant_id")
	tenantID := tenantIDInterface.(uint64)

	koperasiIDStr := c.Query("koperasi_id")
	var koperasiID *uint64

	if koperasiIDStr != "" {
		id, err := strconv.ParseUint(koperasiIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID"})
			return
		}
		koperasiID = &id
	}

	sequences, err := h.sequenceService.GetSequenceList(tenantID, koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sequences": sequences,
	})
}

func (h *SequenceHandler) UpdateSequenceValue(c *gin.Context) {
	tenantIDInterface, _ := c.Get("tenant_id")
	tenantID := tenantIDInterface.(uint64)

	var req UpdateSequenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.sequenceService.UpdateSequenceValue(tenantID, req.KoperasiID, req.SequenceType, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sequence value updated successfully",
	})
}

func (h *SequenceHandler) ResetSequence(c *gin.Context) {
	tenantIDInterface, _ := c.Get("tenant_id")
	tenantID := tenantIDInterface.(uint64)

	var req ResetSequenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.sequenceService.ResetSequence(tenantID, req.KoperasiID, req.SequenceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sequence reset successfully",
	})
}

type UpdateSequenceRequest struct {
	KoperasiID   uint64 `json:"koperasi_id" binding:"required"`
	SequenceType string `json:"sequence_type" binding:"required"`
	Value        uint64 `json:"value" binding:"required"`
}

type ResetSequenceRequest struct {
	KoperasiID   uint64 `json:"koperasi_id" binding:"required"`
	SequenceType string `json:"sequence_type" binding:"required"`
}