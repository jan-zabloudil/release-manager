package repository

import "github.com/nedpals/supabase-go"

type Repository struct{}

func NewRepository(client *supabase.Client) *Repository {
	return &Repository{}
}
