package crypt

import (
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
