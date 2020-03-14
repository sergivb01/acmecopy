package main

import (
	"os"

	"github.com/joho/godotenv"
)

// Config defines the server configuration
type Config struct {
	Listen      string
	PostgresURI string
	Production  bool
	TLSCert     string
	TLSKey      string
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		Listen:      getEnv("LISTEN_ADDR", "localhost:8080"),
		PostgresURI: getEnv("POSTGRES_URI", "postgresql://postgres:sergivb01@127.0.0.1/postgres?sslmode=disable"),
		Production:  getEnv("PRODUCTION", "false") == "true",
		TLSCert:     getEnv("TLS_CERT", "C:\\Users\\Sergi\\Desktop\\acmecopy\\certs\\certificate.pem"),
		TLSKey:      getEnv("TLS_KEY", "C:\\Users\\Sergi\\Desktop\\acmecopy\\certs\\key.pem"),
	}, nil
}

func getEnv(key string, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultVal
}
