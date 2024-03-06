package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateStringRandomly(prefix string, length int) string {
	return prefix + StringWithCharset(length, charset)
}

func StringWithCharset(length int, charset []rune) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func EncryptoThroughMd5(str []byte) string {
	hash := md5.Sum(str)
	hexString := hex.EncodeToString(hash[:])
	return hexString
}
