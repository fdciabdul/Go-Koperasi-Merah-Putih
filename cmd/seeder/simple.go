package main

import (
	"fmt"
	"log"

	"koperasi-merah-putih/config"
	"koperasi-merah-putih/internal/database"
	"koperasi-merah-putih/internal/models/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	dbManager, err := database.NewPostgresConnection(&cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db := dbManager.DB

	fmt.Println("Starting simple seeder...")

	// Create basic tenant
	tenant := &postgres.Tenant{
		TenantCode: "DEMO",
		TenantName: "Demo Tenant",
		Domain:     "demo.local",
		IsActive:   true,
	}
	db.FirstOrCreate(tenant, postgres.Tenant{TenantCode: "DEMO"})
	fmt.Println("✓ Created tenant")

	// Create basic user
	user := &postgres.User{
		TenantID:     tenant.ID,
		Username:     "admin",
		Email:        "admin@demo.local",
		PasswordHash: "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", // password: admin123
		NamaLengkap:  "Administrator",
		Role:         "super_admin",
		IsActive:     true,
	}
	db.FirstOrCreate(user, postgres.User{Email: "admin@demo.local"})
	fmt.Println("✓ Created admin user")

	// Create sequence numbers
	sequences := []postgres.SequenceNumber{
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "global", CurrentNumber: 1},
	}
	for _, seq := range sequences {
		db.FirstOrCreate(&seq, postgres.SequenceNumber{
			TenantID:     seq.TenantID,
			KoperasiID:   seq.KoperasiID,
			SequenceName: seq.SequenceName,
		})
	}
	fmt.Println("✓ Created sequences")

	fmt.Println("Simple seeder completed successfully!")
	fmt.Println("Login credentials:")
	fmt.Println("Email: admin@demo.local")
	fmt.Println("Password: admin123")
}