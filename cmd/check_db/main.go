package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Check schema_migrations table
	rows, err := pool.Query(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		log.Printf("schema_migrations query error: %v", err)
	} else {
		defer rows.Close()
		fmt.Println("Applied migrations:")
		for rows.Next() {
			var v string
			rows.Scan(&v)
			fmt.Println("  ", v)
		}
	}

	// Check merchants table constraints
	rows2, err := pool.Query(ctx, `
		SELECT column_name, is_nullable 
		FROM information_schema.columns 
		WHERE table_name = 'merchants' AND column_name IN ('secret_key', 'public_key')
	`)
	if err != nil {
		log.Printf("columns query error: %v", err)
	} else {
		defer rows2.Close()
		fmt.Println("\nMerchants columns:")
		for rows2.Next() {
			var col, nullable string
			rows2.Scan(&col, &nullable)
			fmt.Printf("  %s: nullable=%s\n", col, nullable)
		}
	}
}