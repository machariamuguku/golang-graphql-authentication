package utils

import (
	"log"

	gonanoid "github.com/matoous/go-nanoid"
)

// RandEmailCodeGenerator : generates a 62 character long random string
func RandEmailCodeGenerator() string {
	// generate a 62 characters long string from a fixed alphabet
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 62)

	// if there's an error generating
	if err != nil {
		// log for the backend
		log.Printf("RandEmailCodeGenerator: error generating random string %v", err)
		return ""
	}

	return id
}
