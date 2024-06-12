package ecoflow

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func encryptHmacSHA256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
