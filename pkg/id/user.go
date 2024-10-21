package id

import "github.com/google/uuid"

type User uuid.UUID

func (u *User) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*u = User(id)
	return nil
}

func (u User) String() string {
	return uuid.UUID(u).String()
}

func (u *User) Scan(data any) error {
	return scanUUID((*uuid.UUID)(u), "User", data)
}

func (u User) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(u).String()), nil
}

func (u *User) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(u), "User", data)
}
