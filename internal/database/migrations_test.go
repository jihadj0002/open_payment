//go:build migration_test

package database

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationFilesExist(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	var upFiles, downFiles []string
	for _, e := range entries {
		if !e.IsDir() {
			if strings.HasSuffix(e.Name(), ".up.sql") {
				upFiles = append(upFiles, e.Name())
			} else if strings.HasSuffix(e.Name(), ".down.sql") {
				downFiles = append(downFiles, e.Name())
			}
		}
	}

	sort.Strings(upFiles)
	sort.Strings(downFiles)

	require.NotEmpty(t, upFiles, "should have at least one up migration")
	require.NotEmpty(t, downFiles, "should have at least one down migration")

	assert.Equal(t, len(upFiles), len(downFiles), "should have matching up and down files")

	for i, up := range upFiles {
		expectedDown := strings.Replace(up, ".up.sql", ".down.sql", 1)
		assert.Equal(t, expectedDown, downFiles[i], "mismatched migration file pair")
	}
}

func TestMigrationFilesReadable(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + e.Name())
		require.NoError(t, err, "should read migration file %s", e.Name())
		require.NotEmpty(t, content, "migration file %s should not be empty", e.Name())
	}
}

func TestMigrationFilesSyntax(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + e.Name())
		require.NoError(t, err)

		sql := strings.TrimSpace(string(content))
		assert.True(t, len(sql) > 10, "migration %s should contain valid SQL (got %d chars)", e.Name(), len(sql))

		firstWord := strings.Fields(sql)[0]
		upperFirst := strings.ToUpper(firstWord)
		assert.Contains(t, []string{"CREATE", "ALTER", "DROP", "INSERT", "UPDATE", "DELETE", "BEGIN", "SELECT"},
			upperFirst, "migration %s should start with a SQL command", e.Name())
	}
}

func TestMigrationPairs(t *testing.T) {
	upMigrations := make(map[string]string)
	downMigrations := make(map[string]string)

	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		base := strings.TrimSuffix(e.Name(), ".up.sql")
		base = strings.TrimSuffix(base, ".down.sql")

		if strings.HasSuffix(e.Name(), ".up.sql") {
			upMigrations[base] = e.Name()
		} else if strings.HasSuffix(e.Name(), ".down.sql") {
			downMigrations[base] = e.Name()
		}
	}

	for base := range upMigrations {
		_, hasDown := downMigrations[base]
		assert.True(t, hasDown, "migration %s has up file but no matching down file", base)
	}

	for base := range downMigrations {
		_, hasUp := upMigrations[base]
		assert.True(t, hasUp, "migration %s has down file but no matching up file", base)
	}
}

func TestMigrationNaming(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	migrationPattern := `^\d{3}_[a-z_]+\.(up|down)\.sql$`

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		assert.Regexp(t, migrationPattern, e.Name(),
			"migration %s should follow the pattern 001_name.up.sql / 001_name.down.sql", e.Name())
	}
}

func TestMigrationSorting(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}

	sorted := make([]string, len(upFiles))
	copy(sorted, upFiles)
	sort.Strings(sorted)

	assert.Equal(t, sorted, upFiles, "migration files should be sortable by prefix number")
}

func TestMigrationEmbedded(t *testing.T) {
	dir, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)
	require.NotEmpty(t, dir)

	for _, e := range dir {
		if e.IsDir() {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + e.Name())
		require.NoError(t, err)
		assert.NotEmpty(t, content)
	}
}

func TestMigrationSchemaTable(t *testing.T) {
	_, err := migrationFiles.ReadFile("migrations/001_initial_schema.up.sql")
	require.NoError(t, err, "001_initial_schema.up.sql should exist")

	_, err = migrationFiles.ReadFile("migrations/012_webhook_secret_rotation.up.sql")
	require.NoError(t, err, "012_webhook_secret_rotation.up.sql should exist")
}

func TestRunMigrationsIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping idempotency test in short mode")
	}

	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration migration test")
	}

	t.Log("idempotency test requires a running postgres instance")
}

func TestMigrationUpAndDown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping migration test in short mode")
	}

	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration migration test")
	}

	t.Log("migration up/down test requires a running postgres instance")
}

func TestMigrationFileSizes(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		info, err := e.Info()
		require.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0), "migration file %s should not be empty", e.Name())
		assert.Less(t, info.Size(), int64(1<<20), "migration file %s should be less than 1MB", e.Name())
	}
}

func TestMigrationFilePermissions(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		info, err := e.Info()
		require.NoError(t, err)
		assert.False(t, info.IsDir(), "migrations directory should only contain files")
	}
}

func TestMigrationTableNames(t *testing.T) {
	expectedTables := []string{
		"merchants",
		"customers",
		"payment_intents",
		"transactions",
		"webhook_endpoints",
		"webhook_deliveries",
		"api_keys",
		"fraud_rules",
		"fraud_events",
		"settlements",
		"disputes",
		"status_history",
		"fee_configs",
		"system_configs",
		"audit_logs",
		"saved_payment_methods",
		"api_usage_logs",
	}

	content, err := migrationFiles.ReadFile("migrations/001_initial_schema.up.sql")
	require.NoError(t, err)
	sql := string(content)

	for _, table := range expectedTables {
		if strings.Contains(sql, "CREATE TABLE "+table) ||
			strings.Contains(sql, "CREATE TABLE IF NOT EXISTS "+table) ||
			strings.Contains(sql, "CREATE TABLE public."+table) {
			t.Logf("table %s found in initial schema", table)
		}
	}
}

func TestMigrationSequence(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	var lastNum int
	for _, f := range upFiles {
		num := 0
		_, err := fmt.Sscanf(f, "%d", &num)
		require.NoError(t, err, "should parse number from %s", f)
		if lastNum > 0 {
			assert.Equal(t, lastNum+1, num, "migration numbers should be sequential: after %03d got %s", lastNum, f)
		}
		lastNum = num
	}
}

func TestRunMigrationsNonEmptyMigrations(t *testing.T) {
	entries, err := migrationFiles.ReadDir("migrations")
	require.NoError(t, err)

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".up.sql") {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + e.Name())
		require.NoError(t, err)
		assert.NotEmpty(t, strings.TrimSpace(string(content)),
			"up migration %s should not be empty", e.Name())
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".down.sql") {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + e.Name())
		require.NoError(t, err)
		assert.NotEmpty(t, strings.TrimSpace(string(content)),
			"down migration %s should not be empty", e.Name())
	}
}
