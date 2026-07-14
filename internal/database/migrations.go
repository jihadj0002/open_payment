package database

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.up.sql migrations/*.down.sql
var migrationFiles embed.FS

func RunMigrations(db *PostgresDB) error {
	_, err := db.Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("reading embedded migrations: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, f := range upFiles {
		var applied bool
		err := db.Pool.QueryRow(
			context.Background(),
			"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)",
			f,
		).Scan(&applied)
		if err != nil {
			return fmt.Errorf("checking migration %s: %w", f, err)
		}
		if applied {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + f)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", f, err)
		}

		tx, err := db.Pool.Begin(context.Background())
		if err != nil {
			return fmt.Errorf("beginning transaction for %s: %w", f, err)
		}

		if _, err := tx.Exec(context.Background(), string(content)); err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("executing migration %s: %w", f, err)
		}

		if _, err := tx.Exec(
			context.Background(),
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			f,
		); err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("recording migration %s: %w", f, err)
		}

		if err := tx.Commit(context.Background()); err != nil {
			return fmt.Errorf("committing migration %s: %w", f, err)
		}

		log.Info().Str("migration", f).Msg("migration applied")
	}

	return nil
}