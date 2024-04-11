package utils

import (
	"math/rand"
	"strings"
	"time"
)

func SecretGenerator(l int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("@#$?%" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < l; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
