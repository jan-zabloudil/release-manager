package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

func ToSvcUser(id, email string, roleString, name, picture any, createdAt, updatedAt time.Time) (svcmodel.User, error) {
	var u svcmodel.User

	if id, err := uuid.Parse(id); err != nil {
		return svcmodel.User{}, err
	} else {
		u.ID = id
	}
	if roleString, ok := roleString.(string); ok {
		role, err := svcmodel.NewUserRole(roleString)
		if err != nil {
			return svcmodel.User{}, err
		}

		u.Role = role
	} else {
		u.Role = svcmodel.NewBasicUserRole()
	}
	if name, ok := name.(string); ok {
		u.Name = name
	}
	if avatarUrl, ok := picture.(string); ok {
		u.AvatarUrl = avatarUrl
	}

	u.Email, u.CreatedAt, u.UpdatedAt = email, createdAt, updatedAt
	return u, nil
}

func ToSvcUsers(dbUsers []supabase.AdminUser) ([]svcmodel.User, error) {
	users := make([]svcmodel.User, 0, len(dbUsers))

	for _, dbUser := range dbUsers {
		user, err := ToSvcUser(
			dbUser.ID,
			dbUser.Email,
			dbUser.AppMetaData["role"],
			dbUser.UserMetaData["name"],
			dbUser.UserMetaData["picture"],
			dbUser.CreatedAt,
			dbUser.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
