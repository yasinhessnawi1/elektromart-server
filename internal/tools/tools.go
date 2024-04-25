package tools

import "github.com/google/uuid"

// GenerateUUID generates a new UUID and returns it as a string.
func GenerateUUID() uint32 {
	newUUID := uuid.New().ID()
	return newUUID
}

func ConvertStringToUint(id string) uint32 {
	newID, err := uuid.Parse(id)
	if err != nil {
		return 0
	}
	return newID.ID()

}
