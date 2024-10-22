package id

import "github.com/google/uuid"

type Deployment uuid.UUID

func NewDeployment() Deployment {
	return Deployment(uuid.New())
}

func (d *Deployment) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*d = Deployment(id)
	return nil
}

func (d Deployment) String() string {
	return uuid.UUID(d).String()
}

func (d *Deployment) Scan(data any) error {
	return scanUUID((*uuid.UUID)(d), "Deployment", data)
}

func (d Deployment) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(d).String()), nil
}

func (d *Deployment) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(d), "Deployment", data)
}
