package id

import "github.com/google/uuid"

type Release uuid.UUID

func NewRelease() Release {
	return Release(uuid.New())
}

func (r Release) IsNil() bool {
	return uuid.UUID(r) == uuid.Nil
}

func (r *Release) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*r = Release(id)
	return nil
}

func (r Release) String() string {
	return uuid.UUID(r).String()
}

func (r *Release) Scan(data any) error {
	return scanUUID((*uuid.UUID)(r), "Release", data)
}

func (r Release) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(r).String()), nil
}

func (r *Release) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(r), "Release", data)
}
