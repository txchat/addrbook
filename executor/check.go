package executor

import (
	"github.com/txchat/addrbook/types"
	"regexp"
)

var baseInfoField = map[string]bool{
	"nickname": true,
	"avatar":   true,
	"phone":    true,
	"email":    true,
	"pubKey":   true,
}

func checkChainNameField(fd string) bool {
	reg := regexp.MustCompile(`chain\.\w+$`)
	return reg.MatchString(fd)
}

func CheckFieldName(fd string) bool {
	v, ok := baseInfoField[fd]
	if ok {
		return v && ok
	}
	return checkChainNameField(fd)
}

func CheckLevel(lv string) bool {
	switch lv {
	case types.LvlPublic:
		fallthrough
	case types.LvlPrivate:
		fallthrough
	case types.LvlProtect:
		return true
	default:
		return false
	}
}
