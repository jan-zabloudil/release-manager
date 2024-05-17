package repository

import (
	"context"
	"errors"
	"net/http"

	"release-manager/pkg/apierrors"
	"release-manager/repository/model"
	"release-manager/repository/query"
	svcmodel "release-manager/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"
)

type UserRepository struct {
	client *supabase.Client
	dbpool *pgxpool.Pool
}

func NewUserRepository(client *supabase.Client, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		client: client,
		dbpool: pool,
	}
}

func (r *UserRepository) Read(ctx context.Context, userID uuid.UUID) (svcmodel.User, error) {
	var u model.User

	err := pgxscan.Get(ctx, r.dbpool, &u, query.ReadUser, pgx.NamedArgs{"id": userID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.User{}, apierrors.NewUserNotFoundError().Wrap(err)
		}

		return svcmodel.User{}, err
	}

	return model.ToSvcUser(u), nil
}

func (r *UserRepository) ReadByEmail(ctx context.Context, email string) (svcmodel.User, error) {
	var u model.User

	err := pgxscan.Get(ctx, r.dbpool, &u, query.ReadUserByEmail, pgx.NamedArgs{"email": email})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.User{}, apierrors.NewUserNotFoundError().Wrap(err)
		}

		return svcmodel.User{}, err
	}

	return model.ToSvcUser(u), nil
}

func (r *UserRepository) ListAll(ctx context.Context) ([]svcmodel.User, error) {
	var u []model.User

	err := pgxscan.Select(ctx, r.dbpool, &u, query.ListUsers)
	if err != nil {
		return nil, err
	}

	return model.ToSvcUsers(u), nil
}

// Delete deletes a user from both the authentication table (auth.users) and the public.users table
// This action must be done via Supabase client
// Because auth.users cannot be accessed directly in the database
// public.users reference the auth.users table, so deleting a user from auth.users will also delete the user from public.users
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.Admin.DeleteUser(ctx, id.String())
	if err != nil {
		var errResponse *supabase.ErrorResponse
		if errors.As(err, &errResponse) && errResponse.Code == http.StatusNotFound {
			return apierrors.NewUserNotFoundError().Wrap(err)
		}

		return err
	}

	return nil
}
