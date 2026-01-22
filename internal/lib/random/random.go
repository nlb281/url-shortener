package random

import (
	"crypto/rand"
	"math/big"
)

func RandomAlias(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var alias []byte

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "Invalid generation", err
		}
		alias = append(alias, chars[num.Int64()])
	}

	return string(alias), nil
}
