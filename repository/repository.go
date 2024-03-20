package repository

import "github.com/nedpals/supabase-go"

type Repository struct {
	User              *UserRepository
	Project           *ProjectRepository
	ProjectInvitation *ProjectInvitationRepository
	ProjectMember     *ProjectMemberRepository
	App               *AppRepository
}

func NewRepository(client *supabase.Client) *Repository {
	return &Repository{
		User:              NewUserRepository(client),
		Project:           NewProjectRepository(client),
		ProjectInvitation: NewProjectInvitationRepository(client),
		ProjectMember:     NewProjectMemberRepository(client),
		App:               NewAppRepository(client),
	}
}
