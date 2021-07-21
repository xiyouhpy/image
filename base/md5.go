package base

import (
	"crypto/md5"
	"fmt"
)

// signMd5 md5加密 sign
const signMd5 = "test"

// GetMd5 进行 md5 加密
func GetMd5(strMd5 string) string {
	if strMd5 == "" {
		return ""
	}

	strTmp := strMd5 + signMd5

	h := md5.New()
	data := []byte(strTmp)

	return fmt.Sprintf("%x", h.Sum(data))
}
