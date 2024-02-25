package repository

import "github.com/nedpals/supabase-go"

type Repository struct {
	User    *UserRepository
	Project *ProjectRepository
}

func NewRepository(client *supabase.Client) *Repository {
	return &Repository{
		User:    &UserRepository{client},
		Project: NewProjectRepository(client),
	}
}
