package tools

import "github.com/google/uuid"

// GenerateUUID generates a new UUID and returns it as a string.
func GenerateUUID() uint32 {
	newUUID := uuid.New().ID()
	return newUUID
}
