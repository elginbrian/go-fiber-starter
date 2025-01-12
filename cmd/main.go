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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = config.MigrateDatabase(db)
	if err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	jwtSecret, err := generateRandomSecret(32)
	if err != nil {
		log.Fatalf("Failed to generate JWT secret: %v", err)
	}

	container := di.NewContainer(db, jwtSecret)

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	routes.SetupRoutes(app, container.UserHandler, container.AuthHandler)

	err = app.Listen(cfg.Port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	fmt.Println("Server started successfully!")
}

func generateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}