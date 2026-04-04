package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	AppConfig *Config
)

type Config struct {
	AppPort          string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTExpireMinutes string
	JWTRefreshToken  string
	JWTExpire        string
	AppUrl           string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}

	AppConfig = &Config{
		AppPort:         getEnv("PORT", "8001"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "hadeboardbe"),
		JWTSecret:       getEnv("JWT_SECRET", "secret"),
		JWTRefreshToken: getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
		JWTExpire:       getEnv("JWT_EXPIRED", "2h"),
		AppUrl:          getEnv("APP_URL", "http://localhost:4000"),
	}
}

func getEnv(key string, fallback string) string {
	value, exist := os.LookupEnv(key)

	if exist {
		return value
	} else {
		return fallback
	}
}

func ConnectDB() {
	cfg := AppConfig
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed connect to database", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
}
