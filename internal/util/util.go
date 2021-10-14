package util

import "log"

// This is here because of nestif linter.
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
