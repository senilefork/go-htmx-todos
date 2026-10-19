package cfg

import "os"

type Config struct {
	Host string
	Port string
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
	return Config{
		Host: host,
		Port: port,
	}
}
