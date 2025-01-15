package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
)

func GetServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
	log.Println("Generating a random JWT secret.")
	return generateRandomSecret(32) 
}

func generateRandomSecret(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Failed to generate random secret: %v", err)
	}
	return hex.EncodeToString(bytes)
}