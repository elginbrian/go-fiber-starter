package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
)

func RunSQLMigrations(db *pgxpool.Pool) error {
	migrationsDir := "../db/migrations"

	_, err := os.Stat(migrationsDir)
	if os.IsNotExist(err) {
		fmt.Println("Migration directory does not exist. Skipping migrations.")
		return nil
	} else if err != nil {
		return fmt.Errorf("could not access migrations directory: %w", err)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("could not read migrations directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		migrationFilePath := filepath.Join(migrationsDir, file.Name())
		fmt.Printf("Running migration: %s\n", migrationFilePath)

		sqlBytes, err := os.ReadFile(migrationFilePath)
		if err != nil {
			return fmt.Errorf("could not read migration file %s: %w", file.Name(), err)
		}

		_, err = db.Exec(context.Background(), string(sqlBytes))
		if err != nil {
			return fmt.Errorf("could not execute migration %s: %w", file.Name(), err)
		}

		fmt.Printf("Migration %s applied successfully!\n", file.Name())
	}

	return nil
}
