package service

import (
	"context"
	"errors"
	"testing"

	githubmock "release-manager/github/mocks"
	repomock "release-manager/repository/mocks"
	svcerr "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSCMRepoService_GetTags(t *testing.T) {
	appID := uuid.New()
	scmRepo, _ := model.NewSCMRepo(appID, "github", "https://github.com/owner/repo")

	testCases := []struct {
		name           string
		mockRepoReturn model.SCMRepo
		mockRepoErr    error
		mockGitReturn  []model.GitTag
		mockGitErr     error
		expectedErr    error
	}{
		{
			name:           "Successfully get tags",
			mockRepoReturn: scmRepo,
			mockRepoErr:    nil,
			mockGitReturn:  []model.GitTag{{Name: "v1.0.0"}},
			mockGitErr:     nil,
			expectedErr:    nil,
		},
		{
			name:           "Error getting repo",
			mockRepoReturn: model.NewEmptySCMRepo(),
			mockRepoErr:    errors.New("unexpected error"),
			expectedErr:    errors.New("unexpected error"),
		},
		{
			name:           "Error listing tags",
			mockRepoReturn: scmRepo,
			mockRepoErr:    nil,
			mockGitReturn:  nil,
			mockGitErr:     errors.New("unexpected error"),
			expectedErr:    errors.New("unexpected error"),
		},
		{
			name:           "Repo not set",
			mockRepoReturn: model.NewEmptySCMRepo(),
			mockRepoErr:    nil,
			expectedErr:    svcerr.ErrSCMRepoNotSet,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(repomock.MockAppRepository)
			mockRepo.On("ReadRepo", mock.Anything, appID).Return(tc.mockRepoReturn, tc.mockRepoErr)

			mockGitHub := new(githubmock.MockGitHubService)
			if tc.name != "Error getting repo" && tc.name != "Repo not set" {
				mockGitHub.On("ListTags", mock.Anything, tc.mockRepoReturn.RepoOwnerIdentifier(), tc.mockRepoReturn.RepoIdentifier()).Return(tc.mockGitReturn, tc.mockGitErr)
			}

			service := SCMRepoService{
				repository: mockRepo,
				github:     mockGitHub,
			}

			_, err := service.GetTags(context.Background(), appID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockGitHub.AssertExpectations(t)
		})
	}
}
