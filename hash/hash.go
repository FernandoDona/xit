package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetHashCode(content []byte) (string, error) {
	h := sha1.New()
	if _, err := h.Write(content); err != nil {
		return "", err
	}

	hashcode := h.Sum(nil)
	return hex.EncodeToString(hashcode), nil
}
