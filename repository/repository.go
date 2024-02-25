package repository

import "github.com/nedpals/supabase-go"

type Repository struct {
	User     *UserRepository
	Template *TemplateRepository
}

func NewRepository(client *supabase.Client) *Repository {
	return &Repository{
		User:     &UserRepository{client},
		Template: NewTemplateRepository(client),
	}
}
