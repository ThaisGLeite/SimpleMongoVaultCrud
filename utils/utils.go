package utils

import (
	"log"
	"os"
	"regexp"
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

func IsStrongPassword(password string) bool {
	// At least 8 characters long, one uppercase letter, one lowercase letter, one number and one special character
	regex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$`
	match, _ := regexp.MatchString(regex, password)
	return match
}
