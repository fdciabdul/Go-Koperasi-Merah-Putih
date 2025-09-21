package handlers

import (
	"net/http"
	"strconv"
	"time"

	"koperasi-merah-putih/internal/repository"
	"koperasi-merah-putih/internal/services"

	"github.com/gin-gonic/gin"
)

type ProdukHandler struct {
	produkService *services.ProdukService
}

func NewProdukHandler(produkService *services.ProdukService) *ProdukHandler {
	return &ProdukHandler{produkService: produkService}
}

// Kategori Produk Handlers
func (h *ProdukHandler) CreateKategoriProduk(c *gin.Context) {
	var req services.CreateKategoriProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.produkService.CreateKategoriProduk(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori produk berhasil dibuat",
		"data":    result,
	})
}

func (h *ProdukHandler) GetKategoriProdukByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.produkService.GetKategoriProdukByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *ProdukHandler) GetAllKategoriProduk(c *gin.Context) {
	result, err := h.produkService.GetAllKategoriProduk()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Satuan Produk Handlers
func (h *ProdukHandler) CreateSatuanProduk(c *gin.Context) {
	var req services.CreateSatuanProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.produkService.CreateSatuanProduk(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Satuan produk berhasil dibuat",
		"data":    result,
	})
}

func (h *ProdukHandler) GetAllSatuanProduk(c *gin.Context) {
	result, err := h.produkService.GetAllSatuanProduk()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Supplier Handlers
func (h *ProdukHandler) CreateSupplier(c *gin.Context) {
	var req services.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	req.CreatedBy = userID.(uint64)

	result, err := h.produkService.CreateSupplier(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Supplier berhasil dibuat",
		"data":    result,
	})
}

func (h *ProdukHandler) GetSuppliersByKoperasi(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
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

	result, err := h.produkService.GetSuppliersByKoperasi(koperasiID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"page":  page,
		"limit": limit,
	})
}

// Produk Handlers
func (h *ProdukHandler) CreateProduk(c *gin.Context) {
	var req services.CreateProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	req.CreatedBy = userID.(uint64)

	result, err := h.produkService.CreateProduk(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk berhasil dibuat",
		"data":    result,
	})
}

func (h *ProdukHandler) GetProdukByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.produkService.GetProdukByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *ProdukHandler) GetProdukByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")
	if barcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Barcode tidak boleh kosong"})
		return
	}

	result, err := h.produkService.GetProdukByBarcode(barcode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *ProdukHandler) GetProduksByKoperasi(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	kategoriIDStr := c.Query("kategori_id")
	nama := c.Query("nama")
	brand := c.Query("brand")
	stokRendahStr := c.Query("stok_rendah")
	readyStockStr := c.Query("ready_stock")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	filters := repository.ProdukFilters{
		Nama:  nama,
		Brand: brand,
	}

	if kategoriIDStr != "" {
		kategoriID, err := strconv.ParseUint(kategoriIDStr, 10, 64)
		if err == nil {
			filters.KategoriID = kategoriID
		}
	}

	if stokRendahStr == "true" {
		filters.StokRendah = true
	}

	if readyStockStr == "true" {
		filters.ReadyStock = true
	}

	result, err := h.produkService.GetProduksByKoperasi(koperasiID, filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"page":  page,
		"limit": limit,
	})
}

func (h *ProdukHandler) GenerateBarcode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	barcode, err := h.produkService.GenerateBarcode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Barcode berhasil dibuat",
		"barcode": barcode,
	})
}

// Purchase Order Handlers
func (h *ProdukHandler) CreatePurchaseOrder(c *gin.Context) {
	var req services.CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	req.CreatedBy = userID.(uint64)

	result, err := h.produkService.CreatePurchaseOrder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Purchase order berhasil dibuat",
		"data":    result,
	})
}

func (h *ProdukHandler) GetPurchaseOrdersByKoperasi(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
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

	result, err := h.produkService.GetPurchaseOrdersByKoperasi(koperasiID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"page":  page,
		"limit": limit,
	})
}

// Pembelian Handlers
func (h *ProdukHandler) CreatePembelian(c *gin.Context) {
	var req services.CreatePembelianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	req.CreatedBy = userID.(uint64)

	result, err := h.produkService.CreatePembelian(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pembelian berhasil dicatat",
		"data":    result,
	})
}

// Penjualan Handlers
func (h *ProdukHandler) CreatePenjualan(c *gin.Context) {
	var req services.CreatePenjualanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	req.CreatedBy = userID.(uint64)

	if req.TanggalTransaksi.IsZero() {
		req.TanggalTransaksi = time.Now()
	}

	result, err := h.produkService.CreatePenjualan(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Penjualan berhasil dicatat",
		"data":    result,
	})
}

// Report Handlers
func (h *ProdukHandler) GetStokReport(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	result, err := h.produkService.GetStokReport(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Laporan stok berhasil diambil",
		"data":    result,
	})
}

func (h *ProdukHandler) GetProdukStokRendah(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	result, err := h.produkService.GetProdukStokRendah(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk dengan stok rendah berhasil diambil",
		"data":    result,
	})
}

func (h *ProdukHandler) GetProdukExpiringSoon(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 30
	}

	result, err := h.produkService.GetProdukExpiringSoon(koperasiID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk yang akan segera expired berhasil diambil",
		"data":    result,
		"days":    days,
	})
}