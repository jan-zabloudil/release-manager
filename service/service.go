package service

import "release-manager/service/model"

type Service struct {
	User *UserService
}

func NewService(ur model.UserRepository) *Service {
	return &Service{
		User: &UserService{ur},
	}
}
