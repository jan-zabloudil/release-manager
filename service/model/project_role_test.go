package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectRole_NewProjectRole(t *testing.T) {
	testCases := []struct {
		name      string
		role      string
		expectErr bool
	}{
		{"Valid admin role", adminProjectRole, false},
		{"Valid editor role", editorProjectRole, false},
		{"Valid viewer role", viewerProjectRole, false},
		{"Invalid role", "invalid_role", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			role, err := NewProjectRole(tc.role)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.role, role.String())
			}
		})
	}
}

func TestProjectRole_IsEqualOrSuperiorTo(t *testing.T) {
	admin, _ := NewProjectRole(adminProjectRole)
	editor, _ := NewProjectRole(editorProjectRole)
	viewer, _ := NewProjectRole(viewerProjectRole)

	testCases := []struct {
		name     string
		role     ProjectRole
		other    ProjectRole
		expected bool
	}{
		{"Admin to Editor", admin, editor, true},
		{"Admin to Viewer", admin, viewer, true},
		{"Editor to Viewer", editor, viewer, true},
		{"Editor to Admin", editor, admin, false},
		{"Viewer to Editor", viewer, editor, false},
		{"Viewer to Admin", viewer, admin, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.role.IsEqualOrSuperiorTo(tc.other)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestRole_IsSuperiorTo(t *testing.T) {
	admin, _ := NewProjectRole(adminProjectRole)
	editor, _ := NewProjectRole(editorProjectRole)
	viewer, _ := NewProjectRole(viewerProjectRole)

	testCases := []struct {
		name     string
		role     ProjectRole
		other    ProjectRole
		expected bool
	}{
		{"Admin to Editor", admin, editor, true},
		{"Admin to Viewer", admin, viewer, true},
		{"Editor to Viewer", editor, viewer, true},
		{"Editor to Admin", editor, admin, false},
		{"Viewer to Editor", viewer, editor, false},
		{"Viewer to Admin", viewer, admin, false},
		{"Admin to Admin", admin, admin, false},
		{"Editor to Editor", editor, editor, false},
		{"Viewer to Viewer", viewer, viewer, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.role.IsSuperiorTo(tc.other)
			assert.Equal(t, tc.expected, res)
		})
	}
}
