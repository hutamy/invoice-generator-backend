package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/hutamy/invoice-generator-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`

	JwtSecret string `env:"JWT_SECRET"`

	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresDB       string `env:"POSTGRES_DB"`
}

var (
	configuration Config
)

func LoadEnv() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	if err := env.Parse(&configuration); err != nil {
		log.Fatalf("failed to parse environment variables: %v", err)
	}

	return configuration
}

func GetConfig() Config {
	return configuration
}

func InitDB(dbUrl string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	migrate(db)
	return db
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.Invoice{},
		&models.InvoiceItem{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
