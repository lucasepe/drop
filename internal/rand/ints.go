package rand

import (
	"crypto/rand"
	"math/big"
)

func Int(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64()) + min
}
