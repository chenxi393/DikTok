package util

import (
	"github.com/bytedance/sonic"
)

func GetLogStr(req interface{}) string {
	str, _ := sonic.MarshalString(req)
	return str
}
