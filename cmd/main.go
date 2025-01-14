package main

import (
	"context"
	"fiber-starter/config"
	"fiber-starter/internal/di"
	"fiber-starter/internal/routes"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	serverPort := config.GetServerPort()
	databaseURL := config.GetDatabaseURL()
	jwtSecret := config.GetJWTSecret()

	db, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Println("Running database migrations...")
	if err := config.MigrateDatabase(); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	container := di.NewContainer(db, jwtSecret)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	routes.SetupRoutes(app, container.UserHandler, container.AuthHandler, container.PostHandler)

	log.Printf("Server is running on port %s", serverPort)
	if err := app.Listen(serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
