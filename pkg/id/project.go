package id

import "github.com/google/uuid"

type Project uuid.UUID

func NewProject() Project {
	return Project(uuid.New())
}

func (p *Project) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*p = Project(id)
	return nil
}

func (p Project) String() string {
	return uuid.UUID(p).String()
}

func (p *Project) Scan(data any) error {
	return scanUUID((*uuid.UUID)(p), "Project", data)
}

func (p Project) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(p).String()), nil
}

func (p *Project) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(p), "Project", data)
}
