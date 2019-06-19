package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type HMAC struct {
	hmac hash.Hash
}

func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))

	return HMAC{hmac: h}
}

func (h HMAC) Hash(input string) (string, error) {
	h.hmac.Reset()
	_, err := h.hmac.Write([]byte(input))
	if err != nil {
		return "", err
	}

	b := h.hmac.Sum(nil)

	return base64.URLEncoding.EncodeToString(b), nil
}
