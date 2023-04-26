package cryptoo

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// MD5 Generate 32-bit MD5 strings
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
