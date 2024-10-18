package id

import "github.com/google/uuid"

type AuthUser uuid.UUID

func NewEmptyAuthUser() AuthUser {
	return AuthUser(uuid.Nil)
}

func (u AuthUser) IsEmpty() bool {
	return uuid.UUID(u) == uuid.Nil
}

func (u *AuthUser) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*u = AuthUser(id)
	return nil
}

func (u AuthUser) String() string {
	return uuid.UUID(u).String()
}

func (u *AuthUser) Scan(data any) error {
	return scanUUID((*uuid.UUID)(u), "AuthUser", data)
}

func (u AuthUser) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(u).String()), nil
}

func (u *AuthUser) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(u), "AuthUser", data)
}
