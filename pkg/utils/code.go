package utils

import (
	"math/rand"
	"time"
)

const (
	// Characters used for verification codes
	codeChars = "0123456789"
	// Length of verification code
	codeLength = 6
)

// GenerateCode generates a random 6-digit verification code
func GenerateCode() string {
	// Create a properly seeded random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		code[i] = codeChars[r.Intn(len(codeChars))]
	}

	return string(code)
}
