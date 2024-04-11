package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func Sha512encode(w string) string {
	sha_512 := sha512.New()
	sha_512.Write([]byte(w))
	return hex.EncodeToString(sha_512.Sum(nil))
}
