package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	postgresUniqueConstraintErrorCode = "23505"
)

func RunTransaction(ctx context.Context, dbpool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	tx, err := dbpool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		err = finishTransaction(ctx, tx, err)
	}()

	err = fn(tx)

	return err
}

// IsUniqueConstraintViolation checks if the error is a unique constraint violation error and that violation happened on the specified constraint
func IsUniqueConstraintViolation(err error, constraintName string) bool {
	var pgConnErr *pgconn.PgError
	if errors.As(err, &pgConnErr) {
		if pgConnErr.Code == postgresUniqueConstraintErrorCode && pgConnErr.ConstraintName == constraintName {
			return true
		}
	}

	return false
}

func finishTransaction(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
			return fmt.Errorf("doing transaction rollback: %w, original error that caused rollback: %w", rollBackErr, err)
		}

		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
