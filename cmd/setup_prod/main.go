package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	// Insert merchants
	merchants := []struct {
		name  string
		email string
		hash  string
	}{
		{"Jihad", "jihad@test.com", "$2a$10$4VyvEe8/1aT8Tijvb/yy3.LPFNUtmRAbUQpiQI6co83FWzwS1AgW2"},
		{"Test User", "test@test.com", "$2a$10$07MzVwxhIgi04l2l/iZxRu048VVAQyOOI0iBsTx03gljrvZTUIOSO"},
	}

	for _, m := range merchants {
		var id string
err := pool.QueryRow(ctx, `
		INSERT INTO merchants (name, email, password_hash, secret_key, public_key, status, created_at, updated_at)
		VALUES ($1, $2, $3, 'sk_live_' || encode(gen_random_bytes(28), 'hex'), 'pk_live_' || encode(gen_random_bytes(28), 'hex'), 'active', NOW(), NOW())
		ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash, status = 'active'
		RETURNING id
	`, m.name, m.email, m.hash).Scan(&id)
		if err != nil {
			log.Printf("Failed to insert %s: %v", m.email, err)
		} else {
			fmt.Printf("Inserted/updated merchant: %s (id: %s)\n", m.email, id)
		}
	}

	// Also create API keys for each
	for _, m := range merchants {
		var merchantID string
		err := pool.QueryRow(ctx, `SELECT id FROM merchants WHERE email = $1`, m.email).Scan(&merchantID)
		if err != nil {
			log.Printf("Failed to get merchant ID for %s: %v", m.email, err)
			continue
		}

		// Create secret API key
		_, err = pool.Exec(ctx, `
			INSERT INTO api_keys (merchant_id, key_prefix, key_hash, name, permissions, status, created_at)
			VALUES ($1, 'sk_live_', encode(sha256(gen_random_bytes(32))::bytea, 'hex'), 'Default Secret Key', '["read","write"]', 'active', NOW())
			ON CONFLICT DO NOTHING
		`, merchantID)
		if err != nil {
			log.Printf("Failed to create secret key for %s: %v", m.email, err)
		}

		// Create publishable API key
		_, err = pool.Exec(ctx, `
			INSERT INTO api_keys (merchant_id, key_prefix, key_hash, name, permissions, status, created_at)
			VALUES ($1, 'pk_live_', encode(sha256(gen_random_bytes(32))::bytea, 'hex'), 'Default Publishable Key', '["read"]', 'active', NOW())
			ON CONFLICT DO NOTHING
		`, merchantID)
		if err != nil {
			log.Printf("Failed to create publishable key for %s: %v", m.email, err)
		}
	}

	fmt.Println("Done!")
}