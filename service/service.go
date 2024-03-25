package service

import "release-manager/service/model"

type Service struct {
	User     *UserService
	Template *TemplateService
}

func NewService(ur model.UserRepository, ttr model.TemplateRepository) *Service {
	return &Service{
		User:     &UserService{ur},
		Template: &TemplateService{ttr},
	}
}
