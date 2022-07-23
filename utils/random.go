package utils

import (
	"math/rand"
	"time"
)

func RandomInt(limit int) int {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	return r.Intn(limit + 1)
}
