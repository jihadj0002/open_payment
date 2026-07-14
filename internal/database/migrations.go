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
		content, err := migrationFiles.ReadFile("migrations/" + f)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", f, err)
		}

		if _, err := db.Pool.Exec(context.Background(), string(content)); err != nil {
			return fmt.Errorf("executing migration %s: %w", f, err)
		}

		log.Info().Str("migration", f).Msg("migration applied")
	}

	return nil
}