package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeUriStr - encode new URI
func EncodeUriStr(data string) string {
	arr := EncodeUri([]byte(data))
	res := hex.EncodeToString(arr)

	return res
}

// EncodeUri - encode new URI
func EncodeUri(data []byte) []byte {
	hasher := md5.New()
	hasher.Write(data)
	md5Sum := hasher.Sum(nil)

	return md5Sum
}
