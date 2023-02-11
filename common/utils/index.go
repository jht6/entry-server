package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetStringFromReader(reader io.ReadCloser) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetServerEnv() string {
	env := os.Getenv("SERVER_ENV")
	// 为安全起见，不设置环境变量时默认按dev环境启动
	if env == "" {
		env = "dev"
	}
	return env
}

func IsProd() bool {
	return GetServerEnv() == "prod"
}

func IsDev() bool {
	env := GetServerEnv()
	return env == "dev" || env == "devcloud"
}

func IsTest() bool {
	env := GetServerEnv()
	return env == "test"
}

func ArrayToString(arr []uint, delim string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(arr), " "), delim), "[]")
}

func StringToMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
