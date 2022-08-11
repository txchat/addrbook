package util

import (
	"encoding/hex"
	"strings"
	"time"

	chattypes "github.com/txchat/addrbook/types"
)

// CheckTimeOut tm 时间，单位毫秒
func CheckTimeOut(tm int64) bool {
	sec := tm / 1000
	nsec := tm % 1000
	t := time.Unix(sec, nsec).Add(chattypes.QueryTimeOut)
	return time.Now().After(t)
}

//HexDecode 兼容0x格式
func HexDecode(in string) ([]byte, error) {
	return hex.DecodeString(strings.Replace(in, "0x", "", 1))
}
