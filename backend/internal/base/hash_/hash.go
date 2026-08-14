package hash_

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"handfree-work/web-restic/base/log_"
)

func Md5(content string) string {
	hash := md5.New()
	hash.Write([]byte(content))
	md5Hash := hash.Sum(nil)
	return hex.EncodeToString(md5Hash)
}

func Base64(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

func Base64Decode(content string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		log_.Error("Base64Decode error: ", err)
		return "", err
	}
	return string(decoded), nil
}

func Base64NoPadding(content string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(content))
}

func Base64DecodeNoPadding(content string) (string, error) {
	decoded, err := base64.RawStdEncoding.DecodeString(content)
	if err != nil {
		log_.Error("Base64Decode error: ", err)
		return "", err
	}
	return string(decoded), nil
}
