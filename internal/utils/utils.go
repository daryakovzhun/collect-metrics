package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func FromPointer[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func ToPointer[T any](p T) *T {
	return &p
}

func ComputeHash(key string, body []byte) ([]byte, error) {
	hash := hmac.New(sha256.New, []byte(key))
	_, err := hash.Write(body)
	if err != nil {
		return nil, fmt.Errorf("failed to write hash: %w", err)
	}

	return hash.Sum(nil), nil
}

func CompareHash(mac1 []byte, mac2 []byte) bool {
	return hmac.Equal(mac1, mac2)
}
