package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func RandString(n int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var b strings.Builder

	for i := 0; i < n; i++ {
		b.WriteByte(charset[rand.Intn(len(charset))])
	}

	return b.String()
}
