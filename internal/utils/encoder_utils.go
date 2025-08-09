package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeURIStr - encode new URI
func EncodeURIStr(data string) string {
	if data == "" {
		return ""
	}
	arr := EncodeURI([]byte(data))
	if arr == nil || len(arr) == 0 {
		return ""
	}
	res := hex.EncodeToString(arr)

	return res
}

// EncodeURI - encode new URI
func EncodeURI(data []byte) []byte {
	if data == nil || len(data) == 0 {
		return nil
	}
	hasher := md5.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil
	}
	md5Sum := hasher.Sum(nil)

	return md5Sum
}
