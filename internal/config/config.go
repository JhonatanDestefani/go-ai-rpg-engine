package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppMode          string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	TelegramBotToken string
	OllamaURL        string
	OllamaModel      string
}

func Load() (*Config, error) {
	_ = godotenv.Overload(".env")

	return &Config{
		AppMode:          getEnv("APP_MODE", "telegram"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "rpg"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		OllamaURL:        getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel:      getEnv("OLLAMA_MODEL", "qwen3:8b"),
	}, nil
}

func (c *Config) DBConnectionString() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
