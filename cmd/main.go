package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/config"
	"koperasi-merah-putih/internal/cache"
	"koperasi-merah-putih/internal/database"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
	cassandraRepo "koperasi-merah-putih/internal/repository/cassandra"
	"koperasi-merah-putih/internal/routes"
	"koperasi-merah-putih/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbManager, err := database.NewDatabaseManager(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to databases: %v", err)
	}
	defer dbManager.Close()

	err = dbManager.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	postgresDB := dbManager.Postgres.DB
	cassandraSession := dbManager.Cassandra.Session

	// Initialize cache (can be nil if not configured)
	var redisCache *cache.RedisCache
	// TODO: Initialize Redis cache if configured

	// Initialize repositories
	userRepo := postgresRepo.NewUserRepository(postgresDB)
	userRegistrationRepo := postgresRepo.NewUserRegistrationRepository(postgresDB)
	koperasiRepo := postgresRepo.NewKoperasiRepository(postgresDB)
	anggotaRepo := postgresRepo.NewAnggotaKoperasiRepository(postgresDB)
	paymentRepo := postgresRepo.NewPaymentRepository(postgresDB)
	paymentProviderRepo := postgresRepo.NewPaymentProviderRepository(postgresDB)
	ppobRepo := postgresRepo.NewPPOBRepository(postgresDB)
	simpanPinjamRepo := postgresRepo.NewSimpanPinjamRepository(postgresDB)
	klinikRepo := postgresRepo.NewKlinikRepository(postgresDB)
	financialRepo := postgresRepo.NewFinancialRepository(postgresDB)
	wilayahRepo := postgresRepo.NewWilayahRepository(postgresDB)
	masterDataRepo := postgresRepo.NewMasterDataRepository(postgresDB)
	sequenceRepo := postgresRepo.NewSequenceRepository(postgresDB)
	produkRepo := postgresRepo.NewProdukRepository(postgresDB)

	// Analytics repository (Cassandra)
	analyticsRepo := cassandraRepo.NewAnalyticsRepository(cassandraSession)

	// Initialize services
	sequenceService := services.NewSequenceService(sequenceRepo)
	paymentService := services.NewPaymentService(paymentRepo, paymentProviderRepo, sequenceService)
	userService := services.NewUserService(userRepo, userRegistrationRepo, anggotaRepo, paymentService, sequenceService)
	ppobService := services.NewPPOBService(ppobRepo, paymentService, sequenceService)
	koperasiService := services.NewKoperasiService(koperasiRepo, anggotaRepo, wilayahRepo, sequenceService)
	simpanPinjamService := services.NewSimpanPinjamService(simpanPinjamRepo, sequenceService)
	klinikService := services.NewKlinikService(klinikRepo, sequenceService)
	financialService := services.NewFinancialService(financialRepo, sequenceService)
	wilayahService := services.NewWilayahService(wilayahRepo)
	masterDataService := services.NewMasterDataService(masterDataRepo)
	produkService := services.NewProdukService(produkRepo, sequenceRepo)
	reportingService := services.NewReportingService(koperasiRepo, anggotaRepo, produkRepo, simpanPinjamRepo, financialRepo, klinikRepo, redisCache)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, userService, ppobService)
	ppobHandler := handlers.NewPPOBHandler(ppobService)
	koperasiHandler := handlers.NewKoperasiHandler(koperasiService)
	simpanPinjamHandler := handlers.NewSimpanPinjamHandler(simpanPinjamService)
	klinikHandler := handlers.NewKlinikHandler(klinikService)
	financialHandler := handlers.NewFinancialHandler(financialService)
	wilayahHandler := handlers.NewWilayahHandler(wilayahService)
	masterDataHandler := handlers.NewMasterDataHandler(masterDataService)
	sequenceHandler := handlers.NewSequenceHandler(sequenceService)
	produkHandler := handlers.NewProdukHandler(produkService)
	reportingHandler := handlers.NewReportingHandler(reportingService)

	// Initialize middleware
	rbacMiddleware := middleware.NewRBACMiddleware(postgresDB)
	auditMiddleware := middleware.NewAuditMiddleware(analyticsRepo)

	// Initialize routes
	appRoutes := routes.NewRoutes(
		userHandler,
		paymentHandler,
		ppobHandler,
		koperasiHandler,
		simpanPinjamHandler,
		klinikHandler,
		financialHandler,
		wilayahHandler,
		masterDataHandler,
		sequenceHandler,
		produkHandler,
		reportingHandler,
		rbacMiddleware,
		auditMiddleware,
	)

	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(middleware.TenantMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"postgres":  "connected",
			"cassandra": "connected",
		})
	})

	appRoutes.SetupRoutes(router)

	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}