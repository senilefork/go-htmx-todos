package cfg

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Host     string
	Port     string
	DbConfig DbConfig
}

type DbConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

// Load reads config from environment with sensible defaults.
func Load() Config {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	return Config{
		Host: host,
		Port: port,
		DbConfig: DbConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Host:     dbHost,
			Port:     os.Getenv("DB_EXTERNAL_PORT"),
		},
	}
}

// DbUrl builds the Postgres DSN. sslmode=disable is fine for a local
// dev container; production would require TLS here.
func (c Config) DbUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DbConfig.User,
		c.DbConfig.Password,
		c.DbConfig.Host,
		c.DbConfig.Port,
		c.DbConfig.Name,
	)
}
