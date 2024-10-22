package id

import "github.com/google/uuid"

type Environment uuid.UUID

func NewEnvironment() Environment {
	return Environment(uuid.New())
}

func (e Environment) IsNil() bool {
	return uuid.UUID(e) == uuid.Nil
}

func (e *Environment) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*e = Environment(id)
	return nil
}

func (e Environment) String() string {
	return uuid.UUID(e).String()
}

func (e *Environment) Scan(data any) error {
	return scanUUID((*uuid.UUID)(e), "Environment", data)
}

func (e Environment) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(e).String()), nil
}

func (e *Environment) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(e), "Environment", data)
}
