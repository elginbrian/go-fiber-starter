package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		GetEnv("DB_USER", "fiberuser"),
		GetEnv("DB_PASSWORD", "fiberpassword"),
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_PORT", "5432"),
		GetEnv("DB_NAME", "fiberdb"),
	)

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	log.Println("Successfully connected to the database.")

	m, err := migrate.New(
		"file://db/migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully.")
}
