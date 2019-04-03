package util

import (
	t_conf "tvm-light/config"
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"crypto/sha256"
)

func StringArrayToByte(strArray []string) [][]byte{
	var args [][]byte;
	for _,v:= range strArray{
		args = append(args,[]byte(v));
	}
	return args;
}

func Sha256(message string) string {
	key := []byte(t_conf.Secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}
