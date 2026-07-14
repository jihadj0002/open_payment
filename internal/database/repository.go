package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepository struct {
	Pool *pgxpool.Pool
}

func NewBaseRepository(db *PostgresDB) *BaseRepository {
	return &BaseRepository{Pool: db.Pool}
}

func (r *BaseRepository) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
