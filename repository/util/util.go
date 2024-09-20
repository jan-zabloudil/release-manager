package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	postgresUniqueConstraintErrorCode = "23505"
)

func FinishTransaction(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("doing transction rollback: %w", rollbackErr)
		}

		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
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
