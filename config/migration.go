package config

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDatabase() error {
	databaseURL := "postgres://user:password@localhost:5432/url_shortener?sslmode=disable"

	migrationDir := "file://db/migrations"

	m, err := migrate.New(
		migrationDir, 
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("could not initialize migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}