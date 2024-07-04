package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"os"

	"golang.org/x/crypto/pbkdf2"
)
func HashSha256(value string) []byte {
	var secretKey string
	secretKey = os.Getenv("SECRET_KEY")
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(value))
	return hash.Sum(nil)
}
func HashSha512(value string) string {
	saltBytes := []byte(os.Getenv("SECRET_KEY"))
	hash := pbkdf2.Key([]byte(value), saltBytes, 1000, 64, sha512.New)
	return hex.EncodeToString(hash)
}
