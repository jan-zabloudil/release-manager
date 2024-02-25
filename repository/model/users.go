package model

import (
	"github.com/google/uuid"
	svcmodel "github.com/jan-zabloudil/release-manager/service/model"
	"github.com/nedpals/supabase-go"
)

func AuthToSvcUser(dbUser supabase.User) (svcmodel.User, error) {
	var svcUser svcmodel.User

	svcUser.Email = dbUser.Email
	svcUser.CreatedAt = dbUser.CreatedAt
	svcUser.UpdatedAt = dbUser.UpdatedAt

	if id, err := uuid.Parse(dbUser.ID); err != nil {
		return svcmodel.User{}, err
	} else {
		svcUser.ID = id
	}
	if name, ok := dbUser.UserMetadata["name"].(string); ok {
		svcUser.Name = name
	}
	if isAdmin, ok := dbUser.AppMetadata["is_admin"].(bool); ok {
		svcUser.IsAdmin = isAdmin
	} else {
		svcUser.IsAdmin = false
	}
	if avatarUrl, ok := dbUser.UserMetadata["picture"].(string); ok {
		svcUser.AvatarUrl = avatarUrl
	}

	return svcUser, nil
}
