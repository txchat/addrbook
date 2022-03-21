package types

import "errors"

// Errors for lottery
var (
	ErrTypeNotExist      = errors.New("类型不存在")
	ErrFieldType         = errors.New("字段类型不合法")
	ErrLevelTypeUnDefine = errors.New("安全级别未定义")
	ErrParams            = errors.New("请求参数错误")
	ErrSignErr           = errors.New("签名不正确")
	ErrLackPermissions   = errors.New("访问权限不足")
	ErrQueryTimeOut      = errors.New("查询请求过期")
	ErrGroupIdFailed     = errors.New("分组id生成失败")
	ErrGroupsArrayOut    = errors.New("分组数量超出限制")
)
