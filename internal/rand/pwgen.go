package rand

import (
	"crypto/rand"
	"math/big"
)

const (
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars  = "0123456789"
	symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?/~"
	allChars    = lowerChars + upperChars + digitChars + symbolChars

	minLen = 8
)

func GenerateSecurePassword(length int) string {
	if length < minLen {
		length = minLen
	}

	password := make([]byte, 0, length)
	password = append(password,
		getRandomChar(lowerChars),
		getRandomChar(upperChars),
		getRandomChar(digitChars),
		getRandomChar(symbolChars),
	)
	for i := 4; i < length; i++ {
		password = append(password, getRandomChar(allChars))
	}

	shuffle(password)

	return string(password)
}

func getRandomChar(charset string) byte {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		panic(err)
	}
	return charset[n.Int64()]
}

func shuffle(data []byte) {
	for i := range data {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(data))))
		if err != nil {
			panic(err)
		}
		data[i], data[j.Int64()] = data[j.Int64()], data[i]
	}
}
