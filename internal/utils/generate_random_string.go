package utils

import (
	"math/rand/v2"
)

func GenerateRandomString(length int) string {
	// Define the characters that can appear in the string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a slice to hold the random string
	result := make([]byte, length)

	// Loop to select random characters
	for i := range result {
		result[i] = charset[rand.IntN(len(charset))]
	}

	// Convert the slice to a string and return
	return string(result)
}
