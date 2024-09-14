package pointer

import "github.com/google/uuid"

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

func UUIDPtr(id uuid.UUID) *uuid.UUID {
	return &id
}
