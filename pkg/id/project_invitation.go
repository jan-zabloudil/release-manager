package id

import "github.com/google/uuid"

type ProjectInvitation uuid.UUID

func NewProjectInvitation() ProjectInvitation {
	return ProjectInvitation(uuid.New())
}

func (i *ProjectInvitation) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*i = ProjectInvitation(id)
	return nil
}

func (i ProjectInvitation) String() string {
	return uuid.UUID(i).String()
}

func (i *ProjectInvitation) Scan(data any) error {
	return scanUUID((*uuid.UUID)(i), "Invitation", data)
}

func (i ProjectInvitation) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(i).String()), nil
}

func (i *ProjectInvitation) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(i), "Invitation", data)
}
