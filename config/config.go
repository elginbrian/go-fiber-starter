package config

import (
	"log"
	"os"
)

func GetServerPort() string {
	port := os.Getenv("PORT") 
	if port == "" {
		port = "3000"
	}
	return port
}

func GetDatabaseURL() string {
	dbURL := os.Getenv("DATABASE_URL") 
	if dbURL == "" {
		log.Fatalf("DATABASE_URL environment variable is not set")
	}
	return dbURL
}

func GetJWTSecret() string {
	jwtSecret := os.Getenv("JWT_SECRET") 
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET environment variable is not set")
	}
	return jwtSecret
}
