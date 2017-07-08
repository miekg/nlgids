package webinvoice

import (
	"math/rand"
	"time"
)

// Kernmerk returns a unique identifier for this invoice.
func Kenmerk(now time.Time) string {
	rand.Seed(time.Now().UnixNano())

	t := now.Format("02012006")
	letters := randString()

	kenmerk := t + letters
	return kenmerk
}

func randString() string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

const length = 2
