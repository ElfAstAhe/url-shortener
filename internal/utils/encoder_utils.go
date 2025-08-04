package utils

import (
	"crypto/sha256"
)

// EncodeUriStr - encode new URI
func EncodeUriStr(data string) []byte {
	return EncodeUri([]byte(data))
}

// EncodeUri - encode new URI
func EncodeUri(data []byte) []byte {
	buf := make([]byte, 32)
	res := sha256.Sum256(data)
	copy(buf, res[:])

	return buf
}
