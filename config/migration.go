package config

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDatabase() error {
	dbConfig := GetDatabaseConfig()

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", 
		dbConfig["user"], 
		dbConfig["password"], 
		dbConfig["host"], 
		dbConfig["port"], 
		dbConfig["dbname"], 
		dbConfig["sslmode"])

	migrationDir := "file://db/migrations" 

	m, err := migrate.New(migrationDir, databaseURL)
	if err != nil {
		return fmt.Errorf("could not initialize migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}