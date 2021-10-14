package env

import (
	"os"
	"strconv"
	"strings"
)

// GetenvString takes an envvar name, gets the value, and returns it if not blank.
// If blank, it will return the provided default value.
func GetenvString(k, defaultValue string) string {
	value := os.Getenv(k)

	if value == "" {
		return defaultValue
	}

	return value
}

// GetenvInt takes an envvar name, gets the value, and returns it if not blank.
// If blank, it will return the provided default value.
func GetenvInt(k string, defaultValue int) int {
	valueString := os.Getenv(k)
	value, err := strconv.Atoi(valueString)

	if err != nil {
		return defaultValue
	}

	return value
}

// GetenvInt64 takes an envvar name, gets the value, and returns it if not blank.
// If blank, it will return the provided default value.
func GetenvInt64(k string, defaultValue int64) int64 {
	valueString := os.Getenv(k)
	value, err := strconv.ParseInt(valueString, 10, 64)

	if err != nil {
		return defaultValue
	}

	return value
}

// GetenvBool takes an envvar name, gets the value, and returns it if not blank.
// If blank, it will return the provided default value.
func GetenvBool(k string, defaultValue bool) bool {
	value := strings.ToLower(os.Getenv(k))
	switch value {
	case "true", "1":
		return true
	case "false", "0":
		return false
	}

	return defaultValue
}
