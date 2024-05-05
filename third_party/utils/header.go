package utils

import (
	"net/http"
	"strconv"
)

func GetHashFromHeader(r *http.Request) string {
	return r.Header.Get("Hash")
}

func GetSizeFromHeader(r *http.Request) int64 {
	size, err := strconv.ParseInt(r.Header.Get("Size"), 10, 64)
	if err != nil {
		return 0
	}
	return size
}
