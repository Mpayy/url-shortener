package util

import (
	"crypto/rand"
	"math/big"
)

const ShortCodeLength = 6
const CharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode() (string, error) {
	result := make([]byte, ShortCodeLength)
	charsetLength := big.NewInt(int64(len(CharSet)))

	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = CharSet[randomIndex.Int64()]
	}

	return string(result), nil
}
