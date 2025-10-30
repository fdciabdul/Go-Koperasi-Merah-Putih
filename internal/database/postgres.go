package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"koperasi-merah-putih/config"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresConnection(cfg *config.PostgresConfig) (*PostgresDB, error) {
	// Use URI format instead of key-value format (fixes Windows connection issues)
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Verify which database we're actually connected to
	var currentDB string
	if err := db.Raw("SELECT current_database()").Scan(&currentDB).Error; err != nil {
		return nil, fmt.Errorf("failed to verify database connection: %v", err)
	}

	if currentDB != cfg.Database {
		return nil, fmt.Errorf("connected to wrong database: expected '%s', got '%s'", cfg.Database, currentDB)
	}

	log.Printf("Connected to PostgreSQL successfully (database: %s)", currentDB)
	return &PostgresDB{DB: db}, nil
}