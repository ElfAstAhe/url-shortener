package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeURIStr - encode new URI
func EncodeURIStr(data string) string {
	arr := EncodeURI([]byte(data))
	res := hex.EncodeToString(arr)

	return res
}

// EncodeURI - encode new URI
func EncodeURI(data []byte) []byte {
	hasher := md5.New()
	hasher.Write(data)
	md5Sum := hasher.Sum(nil)

	return md5Sum
}
