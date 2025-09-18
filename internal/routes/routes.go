package routes

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type Routes struct {
	userHandler         *handlers.UserHandler
	paymentHandler      *handlers.PaymentHandler
	ppobHandler         *handlers.PPOBHandler
	koperasiHandler     *handlers.KoperasiHandler
	simpanPinjamHandler *handlers.SimpanPinjamHandler
	klinikHandler       *handlers.KlinikHandler
	financialHandler    *handlers.FinancialHandler
	wilayahHandler      *handlers.WilayahHandler
	masterDataHandler   *handlers.MasterDataHandler
	sequenceHandler     *handlers.SequenceHandler
	rbacMiddleware      *middleware.RBACMiddleware
	auditMiddleware     *middleware.AuditMiddleware
}

func NewRoutes(
	userHandler *handlers.UserHandler,
	paymentHandler *handlers.PaymentHandler,
	ppobHandler *handlers.PPOBHandler,
	koperasiHandler *handlers.KoperasiHandler,
	simpanPinjamHandler *handlers.SimpanPinjamHandler,
	klinikHandler *handlers.KlinikHandler,
	financialHandler *handlers.FinancialHandler,
	wilayahHandler *handlers.WilayahHandler,
	masterDataHandler *handlers.MasterDataHandler,
	sequenceHandler *handlers.SequenceHandler,
	rbacMiddleware *middleware.RBACMiddleware,
	auditMiddleware *middleware.AuditMiddleware,
) *Routes {
	return &Routes{
		userHandler:         userHandler,
		paymentHandler:      paymentHandler,
		ppobHandler:         ppobHandler,
		koperasiHandler:     koperasiHandler,
		simpanPinjamHandler: simpanPinjamHandler,
		klinikHandler:       klinikHandler,
		financialHandler:    financialHandler,
		wilayahHandler:      wilayahHandler,
		masterDataHandler:   masterDataHandler,
		sequenceHandler:     sequenceHandler,
		rbacMiddleware:      rbacMiddleware,
		auditMiddleware:     auditMiddleware,
	}
}

func (r *Routes) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	api.Use(r.auditMiddleware.TransactionLogger())

	public := api.Group("")
	{
		public.POST("/users/register", r.userHandler.RegisterUser)
		public.POST("/payments/midtrans/callback", r.paymentHandler.HandleMidtransCallback)
		public.POST("/payments/xendit/callback", r.paymentHandler.HandleXenditCallback)
		public.PUT("/users/verify-payment/:payment_id", r.userHandler.VerifyPayment)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PUT("/users/registrations/:id/approve", r.userHandler.ApproveRegistration)
		protected.PUT("/users/registrations/:id/reject", r.userHandler.RejectRegistration)
		protected.POST("/payments", r.paymentHandler.CreatePayment)
	}

	koperasi := api.Group("/koperasi")
	koperasi.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireTenantAccess())
	{
		koperasi.POST("", r.rbacMiddleware.SuperAdminOnly(), r.koperasiHandler.CreateKoperasi)
		koperasi.GET("", r.koperasiHandler.GetKoperasiList)
		koperasi.GET("/:id", r.koperasiHandler.GetKoperasi)
		koperasi.PUT("/:id", r.rbacMiddleware.AdminOnly(), r.koperasiHandler.UpdateKoperasi)
		koperasi.DELETE("/:id", r.rbacMiddleware.SuperAdminOnly(), r.koperasiHandler.DeleteKoperasi)

		koperasi.POST("/anggota", r.rbacMiddleware.AdminOnly(), r.koperasiHandler.CreateAnggota)
		koperasi.GET("/:koperasi_id/anggota", r.koperasiHandler.GetAnggotaList)
		koperasi.GET("/anggota/:id", r.koperasiHandler.GetAnggota)
		koperasi.PUT("/anggota/:id/status", r.rbacMiddleware.AdminOnly(), r.koperasiHandler.UpdateAnggotaStatus)
	}

	wilayah := api.Group("/wilayah")
	{
		wilayah.GET("/provinsi", r.wilayahHandler.GetProvinsiList)
		wilayah.GET("/provinsi/:provinsi_id/kabupaten", r.wilayahHandler.GetKabupatenList)
		wilayah.GET("/kabupaten/:kabupaten_id/kecamatan", r.wilayahHandler.GetKecamatanList)
		wilayah.GET("/kecamatan/:kecamatan_id/kelurahan", r.wilayahHandler.GetKelurahanList)
	}

	simpanPinjam := api.Group("/simpan-pinjam")
	simpanPinjam.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess())
	{
		simpanPinjam.POST("/produk", r.rbacMiddleware.AdminOnly(), r.simpanPinjamHandler.CreateProduk)
		simpanPinjam.GET("/:koperasi_id/produk", r.simpanPinjamHandler.GetProdukList)

		simpanPinjam.POST("/rekening", r.rbacMiddleware.FinancialAccess(), r.simpanPinjamHandler.CreateRekening)
		simpanPinjam.GET("/anggota/:anggota_id/rekening", r.simpanPinjamHandler.GetRekeningByAnggota)

		simpanPinjam.POST("/transaksi", r.rbacMiddleware.FinancialAccess(), r.simpanPinjamHandler.CreateTransaksi)
		simpanPinjam.GET("/rekening/:rekening_id/transaksi", r.simpanPinjamHandler.GetTransaksiByRekening)

		simpanPinjam.GET("/:koperasi_id/statistik", r.rbacMiddleware.AdminOnly(), r.simpanPinjamHandler.GetStatistik)
		simpanPinjam.GET("/pinjaman/jatuh-tempo", r.rbacMiddleware.FinancialAccess(), r.simpanPinjamHandler.GetPinjamanJatuhTempo)
	}

	ppob := api.Group("/ppob")
	{
		ppob.GET("/kategoris", r.ppobHandler.GetKategoriList)
		ppob.GET("/kategoris/:kategori_id/produks", r.ppobHandler.GetProdukByKategori)
		ppob.POST("/transactions", r.rbacMiddleware.PPOBAccess(), r.ppobHandler.CreateTransaction)
	}

	ppobProtected := ppob.Group("")
	ppobProtected.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess())
	{
		ppobProtected.POST("/settlements", r.rbacMiddleware.AdminOnly(), r.ppobHandler.CreateSettlement)
	}

	klinik := api.Group("/klinik")
	klinik.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess(), r.rbacMiddleware.KlinikAccess())
	{
		klinik.POST("/pasien", r.klinikHandler.CreatePasien)
		klinik.GET("/:koperasi_id/pasien", r.klinikHandler.GetPasienList)
		klinik.GET("/pasien/:id", r.klinikHandler.GetPasien)
		klinik.GET("/:koperasi_id/pasien/search", r.klinikHandler.SearchPasien)

		klinik.POST("/tenaga-medis", r.rbacMiddleware.AdminOnly(), r.klinikHandler.CreateTenagaMedis)
		klinik.GET("/:koperasi_id/tenaga-medis", r.klinikHandler.GetTenagaMedisList)

		klinik.POST("/kunjungan", r.klinikHandler.CreateKunjungan)
		klinik.GET("/kunjungan/:id", r.klinikHandler.GetKunjungan)
		klinik.GET("/pasien/:pasien_id/kunjungan", r.klinikHandler.GetKunjunganByPasien)

		klinik.POST("/obat", r.rbacMiddleware.AdminOnly(), r.klinikHandler.CreateObat)
		klinik.GET("/:koperasi_id/obat", r.klinikHandler.GetObatList)
		klinik.GET("/:koperasi_id/obat/search", r.klinikHandler.SearchObat)
		klinik.GET("/:koperasi_id/obat/stok-rendah", r.klinikHandler.GetObatStokRendah)

		klinik.GET("/:koperasi_id/statistik", r.rbacMiddleware.AdminOnly(), r.klinikHandler.GetStatistik)
	}

	financial := api.Group("/financial")
	financial.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess(), r.rbacMiddleware.FinancialAccess())
	{
		financial.POST("/coa/akun", r.rbacMiddleware.AdminOnly(), r.financialHandler.CreateCOAAkun)
		financial.GET("/:koperasi_id/coa/akun", r.financialHandler.GetCOAAkunList)
		financial.GET("/coa/kategori", r.financialHandler.GetCOAKategoriList)

		financial.POST("/jurnal", r.financialHandler.CreateJurnal)
		financial.GET("/:koperasi_id/jurnal", r.financialHandler.GetJurnalList)
		financial.GET("/jurnal/:id", r.financialHandler.GetJurnal)
		financial.PUT("/jurnal/:id/post", r.financialHandler.PostJurnal)
		financial.PUT("/jurnal/:id/cancel", r.financialHandler.CancelJurnal)

		financial.GET("/:koperasi_id/neraca-saldo", r.financialHandler.GetNeracaSaldo)
		financial.GET("/:koperasi_id/laba-rugi", r.financialHandler.GetLabaRugi)
		financial.GET("/:koperasi_id/neraca", r.financialHandler.GetNeraca)
		financial.GET("/akun/:akun_id/saldo", r.financialHandler.GetSaldoAkun)
	}

	masterData := api.Group("/master-data")
	masterData.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireTenantAccess())
	{
		masterData.POST("/kbli", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.CreateKBLI)
		masterData.GET("/kbli", r.masterDataHandler.GetKBLIList)
		masterData.GET("/kbli/:id", r.masterDataHandler.GetKBLI)
		masterData.PUT("/kbli/:id", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.UpdateKBLI)

		masterData.POST("/jenis-koperasi", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.CreateJenisKoperasi)
		masterData.GET("/jenis-koperasi", r.masterDataHandler.GetJenisKoperasiList)
		masterData.GET("/jenis-koperasi/:id", r.masterDataHandler.GetJenisKoperasi)
		masterData.PUT("/jenis-koperasi/:id", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.UpdateJenisKoperasi)

		masterData.POST("/bentuk-koperasi", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.CreateBentukKoperasi)
		masterData.GET("/bentuk-koperasi", r.masterDataHandler.GetBentukKoperasiList)
		masterData.GET("/bentuk-koperasi/:id", r.masterDataHandler.GetBentukKoperasi)
		masterData.PUT("/bentuk-koperasi/:id", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.UpdateBentukKoperasi)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireTenantAccess())
	{
		admin.GET("/sequences", r.rbacMiddleware.AdminOnly(), r.sequenceHandler.GetSequenceList)
		admin.PUT("/sequences/update-value", r.rbacMiddleware.AdminOnly(), r.sequenceHandler.UpdateSequenceValue)
		admin.PUT("/sequences/reset", r.rbacMiddleware.AdminOnly(), r.sequenceHandler.ResetSequence)
	}
}