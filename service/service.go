package service

import "release-manager/service/model"

type Service struct {
	User     *UserService
	Settings *SettingsService
}

func NewService(ur model.UserRepository, sr model.SettingsRepository) *Service {
	return &Service{
		User:     &UserService{ur},
		Settings: &SettingsService{sr},
	}
}
