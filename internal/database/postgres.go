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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode)

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

	log.Println("Connected to PostgreSQL successfully")
	return &PostgresDB{DB: db}, nil
}