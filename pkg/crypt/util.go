package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func MustKeyFromBase64(key string) []byte {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic(fmt.Errorf("failed to decode key: %s", err))
	}

	if len(b) != 32 {
		panic(fmt.Sprintf("invalid CRYPT_KEY length: got %d, want 32 bytes", len(key)))
	}

	return b
}

func GenerateBase64Key() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(key), nil
}
