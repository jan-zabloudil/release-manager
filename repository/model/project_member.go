package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateProjectMemberInput struct {
	UserID      uuid.UUID `json:"user_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	ProjectRole string    `json:"project_role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectMember struct {
	User        User      `json:"users"` // Supabase returns joined table data in json array named after joined table, "users" in this case
	ProjectID   uuid.UUID `json:"project_id"`
	ProjectRole string    `json:"project_role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
