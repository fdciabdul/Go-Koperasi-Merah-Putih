package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/config"
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

	userRepo := postgresRepo.NewUserRepository(postgresDB)
	userRegistrationRepo := postgresRepo.NewUserRegistrationRepository(postgresDB)
	paymentRepo := postgresRepo.NewPaymentRepository(postgresDB)
	paymentProviderRepo := postgresRepo.NewPaymentProviderRepository(postgresDB)
	ppobRepo := postgresRepo.NewPPOBRepository(postgresDB)
	analyticsRepo := cassandraRepo.NewAnalyticsRepository(cassandraSession)

	sequenceService := services.NewSequenceService(postgresDB)
	paymentService := services.NewPaymentService(paymentRepo, paymentProviderRepo, sequenceService)
	userService := services.NewUserService(userRepo, userRegistrationRepo, paymentService, sequenceService)
	ppobService := services.NewPPOBService(ppobRepo, paymentService, sequenceService)

	userHandler := handlers.NewUserHandler(userService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, userService, ppobService)
	ppobHandler := handlers.NewPPOBHandler(ppobService)

	routes := routes.NewRoutes(userHandler, paymentHandler, ppobHandler)

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

	routes.SetupRoutes(router)

	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}