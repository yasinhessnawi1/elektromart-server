package tools

import (
	"github.com/google/uuid"
	"strconv"
)

// GenerateUUID generates a new UUID and returns it as a uint32.
// This function uses the UUID package to create a new unique identifier, which is then converted to uint32.
// This is particularly useful for generating unique IDs for database entries that require numeric identifiers.
func GenerateUUID() uint32 {
	newUUID := uuid.New().ID() // Generate a new UUID and retrieve the numeric representation.
	return newUUID
}

// ConvertStringToUint attempts to convert a string representation of a number into a uint32.
// This is useful for parsing numeric strings from sources like HTTP requests into usable uint32 IDs.
// Returns the uint32 representation if the conversion is successful; otherwise, returns 0.
// This function provides a basic error handling mechanism where it returns 0 if the conversion fails,
// which can be checked by the caller to determine if the conversion was successful.
func ConvertStringToUint(id string) uint32 {
	newID, err := strconv.ParseUint(id, 10, 32) // Attempt to parse the string as a base-10 uint32.
	if err != nil {
		return 0 // Return 0 if there is an error during conversion, indicating failure.
	}
	return uint32(newID) // Return the converted uint32 value.
}
