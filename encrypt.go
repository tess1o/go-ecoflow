package ecoflow

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Step 4: encrypt
// E.g. byte[] signBytes = HMAC-SHA256(str, secretKey)
// Step 5: convert byte[] to hexadecimal string. String sign = bytesToHexString(signBytes)
// E.g. sign=85776ede686fe4783eac48135b0b1748ba2d7e9bb7791b826dc942fc29d4ada8
func encryptHmacSHA256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
