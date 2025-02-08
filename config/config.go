package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetServerPort() string {
	port := ":8080"
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
    return "my-very-secure-secret"
}

func GetRefreshSecret() string {
	return "my-very-secure-refresh-secret"
}

func InitDatabase() *pgxpool.Pool {
	databaseURL := GetDatabaseURL()
	db, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Println("Running database migrations...")
	if err := RunSQLMigrations(db); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	return db
}


func SetupFiber() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	return app
}

func StartServer(app *fiber.App, port string) {
	log.Printf("Server is running on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}