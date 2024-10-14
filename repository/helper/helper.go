package helper

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	postgresUniqueConstraintErrorCode = "23505"
)

type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type ExecExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

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

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func ReadValue[T any](ctx context.Context, q Querier, query string, args pgx.NamedArgs) (T, error) {
	var result T
	if err := pgxscan.Get(ctx, q, &result, query, args); err != nil {
		return result, err
	}
	return result, nil
}

func ListValues[T any](ctx context.Context, q Querier, query string, args pgx.NamedArgs) ([]T, error) {
	var result []T
	if err := pgxscan.Select(ctx, q, &result, query, args); err != nil {
		return result, err
	}
	return result, nil
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
