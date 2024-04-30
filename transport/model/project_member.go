package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMember struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	ProjectRole string    `json:"project_role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
