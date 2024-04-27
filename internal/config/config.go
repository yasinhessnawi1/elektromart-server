package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// LoadConfig loads environment variables from a .env file located in the same directory as the application.
// It uses the godotenv package to load the file. If no .env file is found, it logs a message but does not terminate the application.
// This function is typically called at the start of the application to initialize environment settings.
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found") // Log a message if no .env file is found. This is not a critical error.
	}
}

// GetConfig retrieves the value of an environment variable given its key.
// It returns the value associated with the key as a string.
// If the key does not exist in the environment, it returns an empty string.
// This function is used to abstract the access to environment variables, allowing for easy retrieval of config settings.
func GetConfig(key string) string {
	return os.Getenv(key) // Directly return the value of the environment variable.
}
