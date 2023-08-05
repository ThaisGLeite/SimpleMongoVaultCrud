package utils

import (
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// Regular expressions to match password requirements
var (
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex = regexp.MustCompile(`[a-z]`)
	numberRegex    = regexp.MustCompile(`[0-9]`)
	specialRegex   = regexp.MustCompile(`[!@#$%^&*-]`)
)

// Initialize logger with JSON formatter, standard output and info log level
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// GetEnv retrieves the value of the environment variable named by the key.
// If the environment variable is not set, it returns the provided fallback value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// HandleError logs the error based on the provided log level: "I" for Info, "E" for Error, and "W" for Warning.
// If the error is not nil, it logs the error and the associated message at the given log level.
func HandleError(logLevel, msg string, err error) {
	if err != nil {
		// Create an entry for the error
		entry := log.WithFields(log.Fields{
			"error": err,
		})

		// Log the error at the appropriate level
		switch logLevel {
		case "I":
			entry.Info(msg)
		case "E":
			entry.Error(msg)
		case "W":
			entry.Warn(msg)
		default:
			entry.Info(msg)
		}
	}
}

// IsStrongPassword checks whether the provided password string fulfills all password requirements.
// It checks for minimum length of 8, at least one uppercase letter, one lowercase letter,
// one digit and one special character. It returns true if all requirements are met, false otherwise.
func IsStrongPassword(password string) bool {
	hasMinLength := len(password) >= 8
	hasUppercase := uppercaseRegex.MatchString(password)
	hasLowercase := lowercaseRegex.MatchString(password)
	hasNumber := numberRegex.MatchString(password)
	hasSpecial := specialRegex.MatchString(password)

	return hasMinLength && hasUppercase && hasLowercase && hasNumber && hasSpecial
}
