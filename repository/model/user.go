package model

import (
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/relvacode/iso8601"
)

type User struct {
	ID        uuid.UUID    `json:"id"`
	Email     string       `json:"email"`
	Name      string       `json:"name"`
	AvatarURL string       `json:"avatar_url"`
	Role      string       `json:"role"`
	CreatedAt iso8601.Time `json:"created_at"`
	UpdatedAt iso8601.Time `json:"updated_at"`
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
			dbUser.CreatedAt.Time,
			dbUser.UpdatedAt.Time,
		)
		if err != nil {
			return nil, err
		}

		u = append(u, user)
	}

	return u, nil
}
