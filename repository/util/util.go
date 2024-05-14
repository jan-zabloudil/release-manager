package util

import (
	"context"
	"errors"
	"fmt"

	"release-manager/pkg/dberrors"

	"github.com/jackc/pgx/v5"
	postgrestgo "github.com/nedpals/postgrest-go/pkg"
)

const (
	postgresSingleRecordFetchErrorCode = "PGRST116"
)

func ToDBError(err error) *dberrors.DBError {
	var postgreErr *postgrestgo.RequestError
	if errors.As(err, &postgreErr) {
		if postgreErr.Code == postgresSingleRecordFetchErrorCode {
			return dberrors.NewNotFoundError().Wrap(err)
		}
	}

	return dberrors.NewUnknownError().Wrap(err)
}

func FinishTransaction(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("failed to rollback tx: %w", rollbackErr)
		}

		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
