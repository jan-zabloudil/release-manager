package service

import "github.com/jan-zabloudil/release-manager/service/model"

type Service struct {
	User *UserService
}

func NewService(us model.UserRepository) *Service {
	return &Service{
		User: &UserService{us},
	}
}
