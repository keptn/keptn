package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// State generates a new pseudo random base64 encoded string of length n
// to be used to avoid CSRF attacks. Note that n must be >=1
func State(n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("length must be at least one")
	}
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
