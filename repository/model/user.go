package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToSvcUsers(users []User) ([]svcmodel.User, error) {
	u := make([]svcmodel.User, 0, len(users))

	for _, dbUser := range users {
		user, err := svcmodel.ToUser(
			dbUser.ID,
			dbUser.Email,
			dbUser.Name,
			dbUser.AvatarURL,
			dbUser.Role,
			dbUser.CreatedAt,
			dbUser.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		u = append(u, user)
	}

	return u, nil
}
