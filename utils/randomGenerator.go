package utils

import (
	"log"

	gonanoid "github.com/matoous/go-nanoid"
)

// RandGenerator : generates a 62 character long random string
func RandGenerator() string {
	// generate a 62 characters long string from a fixed alphabet
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 62)

	// if there's an error generating
	if err != nil {
		// log for the backend
		log.Printf("RandGenerator: error generating random string %v", err)
		return ""
	}

	return id
}
