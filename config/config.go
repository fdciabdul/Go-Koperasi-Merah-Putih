package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	// Database
	Postgres  PostgresConfig
	Cassandra CassandraConfig

	// Application
	App AppConfig

	// Payment Gateway
	Payment PaymentConfig

	// PPOB
	PPOB PPOBConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type CassandraConfig struct {
	Hosts       []string
	Keyspace    string
	Username    string
	Password    string
	Consistency string
}

type AppConfig struct {
	Environment string
	Port        string
	SecretKey   string
	JWTSecret   string
	JWTExpire   int
}

type PaymentConfig struct {
	Midtrans MidtransConfig
	Xendit   XenditConfig
}

type MidtransConfig struct {
	ServerKey   string
	ClientKey   string
	MerchantID  string
	Environment string
}

type XenditConfig struct {
	SecretKey    string
	WebhookToken string
}

type PPOBConfig struct {
	ProviderURL string
	APIKey      string
	SecretKey   string
}

func LoadConfig() (*Config, error) {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	viper.AutomaticEnv()

	config := &Config{
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			Database: getEnv("POSTGRES_DB", "koperasi_merah_putih"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Cassandra: CassandraConfig{
			Hosts:       []string{getEnv("CASSANDRA_HOSTS", "127.0.0.1")},
			Keyspace:    getEnv("CASSANDRA_KEYSPACE", "koperasi_analytics"),
			Username:    getEnv("CASSANDRA_USERNAME", "cassandra"),
			Password:    getEnv("CASSANDRA_PASSWORD", "cassandra"),
			Consistency: getEnv("CASSANDRA_CONSISTENCY", "quorum"),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("APP_PORT", "8080"),
			SecretKey:   getEnv("APP_SECRET_KEY", "default-secret-key"),
			JWTSecret:   getEnv("JWT_SECRET", "jwt-secret-key"),
			JWTExpire:   24,
		},
		Payment: PaymentConfig{
			Midtrans: MidtransConfig{
				ServerKey:   getEnv("MIDTRANS_SERVER_KEY", ""),
				ClientKey:   getEnv("MIDTRANS_CLIENT_KEY", ""),
				MerchantID:  getEnv("MIDTRANS_MERCHANT_ID", ""),
				Environment: getEnv("MIDTRANS_ENVIRONMENT", "sandbox"),
			},
			Xendit: XenditConfig{
				SecretKey:    getEnv("XENDIT_SECRET_KEY", ""),
				WebhookToken: getEnv("XENDIT_WEBHOOK_TOKEN", ""),
			},
		},
		PPOB: PPOBConfig{
			ProviderURL: getEnv("PPOB_PROVIDER_URL", ""),
			APIKey:      getEnv("PPOB_API_KEY", ""),
			SecretKey:   getEnv("PPOB_SECRET_KEY", ""),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}