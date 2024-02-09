package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/zerodoctor/zdcli/logger"
)

func RandString(n int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var b strings.Builder

	for i := 0; i < n; i++ {
		b.WriteByte(charset[rand.Intn(len(charset))])
	}

	return b.String()
}

func InArray(source string, strs []string) bool {
	for _, str := range strs {
		if source == str {
			return true
		}
	}

	return false
}

var ErrNilInterface error = errors.New("interface is nil")

func StructString(s interface{}) (string, error) {
	if s == nil {
		logger.Warnf("failed to convert to string [error=%s]", ErrNilInterface)
		return "", nil
	}

	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to marshal response [error=%s]", err.Error())
	}

	return string(b), nil
}
