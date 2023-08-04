package utils

import (
	"log"
	"os"
)

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the value, which will be set to fallback if the variable is not present.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func HandleError(msg string, err error) {
	if err != nil {
		log.Panic(msg, err)
	}
}
