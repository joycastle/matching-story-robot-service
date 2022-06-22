package lib

import (
	"crypto/md5"
	"fmt"
)

func Md5(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
