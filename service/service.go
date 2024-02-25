package service

import "release-manager/service/model"

type Service struct {
	User    *UserService
	Project *ProjectService
}

func NewService(ur model.UserRepository, pr model.ProjectRepository) *Service {
	return &Service{
		User:    &UserService{ur},
		Project: &ProjectService{pr},
	}
}
