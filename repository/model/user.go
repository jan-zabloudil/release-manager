package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

func ToSvcUser(
	id, email string,
	isAdmin, name, picture interface{},
	createdAt, updatedAt time.Time,
) (svcmodel.User, error) {
	var u svcmodel.User

	if id, err := uuid.Parse(id); err != nil {
		return svcmodel.User{}, err
	} else {
		u.ID = id
	}
	if isAdmin, ok := isAdmin.(bool); ok {
		u.IsAdmin = isAdmin
	} else {
		u.IsAdmin = false
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
