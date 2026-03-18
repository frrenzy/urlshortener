package service

import (
	"math/rand"
	"time"
)

const charSet string = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateShort() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charSet[seededRand.Intn(len(charSet)-1)]
	}
	return string(b)
}
