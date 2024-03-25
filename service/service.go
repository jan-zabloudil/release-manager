package service

import "release-manager/service/model"

type Service struct {
	User              *UserService
	Project           *ProjectService
	ProjectInvitation *ProjectInvitationService
	ProjectMember     *ProjectMemberService
	ProjectMembership *ProjectMembershipManagementService
	App               *AppService
	SCMRepo           *SCMRepoService
	Release           *ReleaseService
}

func NewService(
	ur model.UserRepository,
	pr model.ProjectRepository,
	pir model.ProjectInvitationRepository,
	pmr model.ProjectMemberRepository,
	ar model.AppRepository,
	sr model.SCMRepoRepository,
	github model.GitHub,
	rr model.ReleaseRepository,
	slack model.Slack,
) *Service {
	userSvc := &UserService{ur}
	projectSvc := &ProjectService{pr}
	appSvc := &AppService{ar}

	projectInvitationSvc := &ProjectInvitationService{pir}
	projectMemberSvc := &ProjectMemberService{pmr}

	return &Service{
		User:              userSvc,
		Project:           projectSvc,
		ProjectInvitation: projectInvitationSvc,
		ProjectMember:     projectMemberSvc,
		ProjectMembership: &ProjectMembershipManagementService{
			userSvc:       userSvc,
			memberSvc:     projectMemberSvc,
			invitationSvc: projectInvitationSvc,
		},
		App:     appSvc,
		SCMRepo: &SCMRepoService{sr, github},
		Release: &ReleaseService{appSvc, projectSvc, rr, slack},
	}
}
