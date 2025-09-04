package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultDBHost     = "localhost"
	defaultDBUser     = "postgres"
	defaultDBPassword = ""
	defaultDBName     = "sub_service"
	defaultDBPort     = "5432"
	defaultHTTPPort   = "8080"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	ServerPort string
}

func InitConfig(logger *log.Logger) *Config {
	if err := godotenv.Load(); err != nil {
		logger.Println(".env not found, using default variables")
	}
	return &Config{
		DBHost:     getEnv("DB_HOST", defaultDBHost),
		DBUser:     getEnv("DB_USER", defaultDBUser),
		DBPassword: getEnv("DB_PASSWORD", defaultDBPassword),
		DBName:     getEnv("DB_NAME", defaultDBName),
		DBPort:     getEnv("DB_PORT", defaultDBPort),
		ServerPort: getEnv("SERVER_PORT", defaultHTTPPort),
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}
