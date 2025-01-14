package config

import (
	"log"
	"os"
)

func GetServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000" 
	}
	return port
}

func GetDatabaseConfig() map[string]string {
	dbConfig := map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
		"dbname":   os.Getenv("DB_NAME"),
		"sslmode":  os.Getenv("DB_SSLMODE"), 
	}
	
	for key, value := range dbConfig {
		if value == "" && key != "sslmode" {
			log.Fatalf("Missing required environment variable for database: %s", key)
		}
	}
	return dbConfig
}