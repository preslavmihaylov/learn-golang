package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const tokenSize = 32

func RememberToken() (string, error) {
	return String(tokenSize)
}

func String(nBytes int) (string, error) {
	b, err := randBytes(nBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func randBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
