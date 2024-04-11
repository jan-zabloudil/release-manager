package repository

import "github.com/nedpals/supabase-go"

type Repository struct {
	Auth *AuthRepository
	User *UserRepository
}

func NewRepository(client *supabase.Client) *Repository {
	return &Repository{
		Auth: NewAuthRepository(client),
		User: NewUserRepository(client),
	}
}
