package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectMember_HasAtLeastRole(t *testing.T) {
	adminRole, _ := NewProjectRole(adminProjectRole)
	editorRole, _ := NewProjectRole(editorProjectRole)
	viewerRole, _ := NewProjectRole(viewerProjectRole)

	cases := []struct {
		name        string
		memberRole  ProjectRole
		compareRole ProjectRole
		want        bool
	}{
		{"Admin vs Editor", adminRole, editorRole, true},
		{"Editor vs Viewer", editorRole, viewerRole, true},
		{"Viewer vs Admin", viewerRole, adminRole, false},
		{"Editor vs Admin", editorRole, adminRole, false},
		{"Viewer vs Editor", viewerRole, editorRole, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			member := ProjectMember{
				Role: c.memberRole,
			}
			result := member.HasAtLeastRole(c.compareRole)
			assert.Equal(t, c.want, result)
		})
	}
}

func TestProjectMember_CanGrantRole(t *testing.T) {
	adminRole, _ := NewProjectRole(adminProjectRole)
	editorRole, _ := NewProjectRole(editorProjectRole)
	viewerRole, _ := NewProjectRole(viewerProjectRole)

	cases := []struct {
		name       string
		memberRole ProjectRole
		grantRole  ProjectRole
		want       bool
	}{
		{"Admin grants Editor", adminRole, editorRole, true},
		{"Editor grants Viewer", editorRole, viewerRole, true},
		{"Viewer grants Editor", viewerRole, editorRole, false},
		{"Editor grants Admin", editorRole, adminRole, false},
		{"Viewer grants Admin", viewerRole, adminRole, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			member := ProjectMember{
				Role: c.memberRole,
			}
			result := member.CanGrantRole(c.grantRole)
			assert.Equal(t, c.want, result)
		})
	}
}

func TestProjectMember_CanUpdateMember(t *testing.T) {
	adminRole, _ := NewProjectRole(adminProjectRole)
	editorRole, _ := NewProjectRole(editorProjectRole)
	viewerRole, _ := NewProjectRole(viewerProjectRole)

	cases := []struct {
		name        string
		memberRole  ProjectRole
		otherMember ProjectMember
		want        bool
	}{
		{"Admin updates Editor", adminRole, ProjectMember{Role: editorRole}, true},
		{"Editor updates Viewer", editorRole, ProjectMember{Role: viewerRole}, true},
		{"Viewer updates Editor", viewerRole, ProjectMember{Role: editorRole}, false},
		{"Editor updates Admin", editorRole, ProjectMember{Role: adminRole}, false},
		{"Viewer updates Admin", viewerRole, ProjectMember{Role: adminRole}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			member := ProjectMember{
				Role: c.memberRole,
			}
			result := member.CanUpdateMember(c.otherMember)
			assert.Equal(t, c.want, result)
		})
	}
}
