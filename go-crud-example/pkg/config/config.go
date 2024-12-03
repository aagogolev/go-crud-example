package config

import (
    "fmt"
    "github.com/joho/godotenv"
    "log"
    "os"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
}

type ServerConfig struct {
    Port string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func LoadConfig() (*Config, error) {
    envPaths := []string{
        "../.env",
        "./.env", 
        "/app/.env",
    }

    var loadError error
    for _, path := range envPaths {
        if _, err := os.Stat(path); err == nil {
            loadError = godotenv.Load(path)
            if loadError == nil {
                log.Printf("Loaded environment from %s", path)
                break
            }
        }
    }

    config := &Config{
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8000"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            DBName:   getEnv("DB_NAME", "users"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
    }

    return config, nil
}

func (d *DatabaseConfig) GetDSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
    )
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
