package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateAPIKey(length int) (string, error) {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
