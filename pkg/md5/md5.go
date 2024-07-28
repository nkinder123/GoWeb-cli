package md5

import (
	"crypto/md5"
	"encoding/hex"
)

const secret = "bluebell"

// md5加密过程
func EncodeMd5(opassword string) string {
	str := md5.New()
	str.Write([]byte(secret))
	return hex.EncodeToString(str.Sum([]byte(opassword)))
}
