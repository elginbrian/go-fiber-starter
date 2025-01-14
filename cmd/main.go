package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fiber-starter/config"
	"fiber-starter/internal/di"
	"fiber-starter/internal/routes"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.DatabaseURL == "" || cfg.Port == "" {
		log.Fatal("Missing required configuration (DatabaseURL or Port)")
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	if err := config.MigrateDatabase(); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret, err = generateRandomSecret(32)
		if err != nil {
			log.Fatalf("Failed to generate JWT secret: %v", err)
		}
	}

	container := di.NewContainer(db, jwtSecret)

	app := fiber.New()
	app.Use(logger.New()) 
	app.Use(cors.New())

	routes.SetupRoutes(app, container.UserHandler, container.AuthHandler, container.PostHandler)

	fmt.Printf("Server running on port %s\n", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func generateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}