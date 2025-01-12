package config

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDatabase(db *sql.DB) error {
	databaseURL := "postgres://user:password@localhost:5432/url_shortener?sslmode=disable"

	m, err := migrate.New(
		"file://../db/migrations", 
		databaseURL,               
	)
	if err != nil {
		return fmt.Errorf("could not initialize migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}