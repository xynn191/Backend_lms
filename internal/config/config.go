package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	DBUser       string
	DBPass       string
	DBHost       string
	DBPort       string
	DBName       string
	JWTSecret    string
	JWTExpiresIn time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load() // Load .env file, ignore error if not found

	expStr := getEnv("JWT_EXPIRES_IN", "24h")
	exp, err := time.ParseDuration(expStr)
	if err != nil {
		exp = 24 * time.Hour
	}

	return &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPass:       getEnv("DB_PASS", ""),
		DBHost:       getEnv("DB_HOST", "127.0.0.1"),
		DBPort:       getEnv("DB_PORT", "3306"),
		DBName:       getEnv("DB_NAME", "edulms_db"),
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		JWTExpiresIn: exp,
	}, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
