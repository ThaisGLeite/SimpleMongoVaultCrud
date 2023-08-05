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
	hasMinLength := len(password) >= 8
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*-]`).MatchString(password)

	return hasMinLength && hasUppercase && hasLowercase && hasNumber && hasSpecial
}
