package model

import (
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
)

type User struct {
	ID        id.User   `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `json:"user_role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUser(u svcmodel.User) User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUsers(users []svcmodel.User) []User {
	u := make([]User, 0, len(users))
	for _, user := range users {
		u = append(u, ToUser(user))
	}

	return u
}

type AuthUser struct {
	User               User                `json:"user"`
	ProjectMemberships []ProjectMembership `json:"project_memberships"`
}

type ProjectMembership struct {
	ProjectID   id.Project `json:"project_id"`
	ProjectRole string     `json:"project_role"`
}

func ToProjectMemberships(pm []svcmodel.ProjectMember) []ProjectMembership {
	m := make([]ProjectMembership, 0, len(pm))
	for _, member := range pm {
		m = append(m, ProjectMembership{
			ProjectID:   member.ProjectID,
			ProjectRole: string(member.ProjectRole),
		})
	}

	return m
}

func ToAuthUser(u svcmodel.User, m []svcmodel.ProjectMember) AuthUser {
	return AuthUser{
		User:               ToUser(u),
		ProjectMemberships: ToProjectMemberships(m),
	}
}
