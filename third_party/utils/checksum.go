package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func Checksum(r io.Reader) string {
	hash := sha256.New()
	_, _ = io.Copy(hash, r)
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
