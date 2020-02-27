package utils

import (
	"math/rand"
	"time"
)

// generate random string between min - max and add min to the result
func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

// RandPhoneCodeGenerator generates a 6 digit long random integer
func RandPhoneCodeGenerator() int {
	// seed a constantly changing stream
	// like timestring
	rand.Seed(time.Now().UnixNano())
	// generate
	randomNum := random(123456, 987604)

	return randomNum
}
