package uuid

import (
	"github.com/google/uuid"
)

type UUID = uuid.UUID

func UUIDFromBytes(b []byte) (UUID, error) {
	return uuid.FromBytes(b)
}

func UUIDFromString(s string) (UUID, error) {
	return uuid.Parse(s)
}
